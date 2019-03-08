package server

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	myapp "github.com/clementauger/monparcours/server/app"
	"github.com/clementauger/monparcours/server/dbconnect"
	"github.com/clementauger/monparcours/server/env"
	mysqlmodel "github.com/clementauger/monparcours/server/model/mysql"
	// pgsqlmodel "github.com/clementauger/monparcours/server/model/pgsql"
	"github.com/gobuffalo/packr"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/dchest/captcha"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Getkey(ctx context.Context) {
	flag.Parse()

	stage := env.Stage()

	var appConfig myapp.Environment
	{
		env, err := myapp.GetEnvironment("app.yml", stage)
		if err != nil {
			log.Fatal(err)
		}
		appConfig = *env
	}

	fmt.Print(myapp.GetKey(appConfig))
}

func ServeHTTP(ctx context.Context) {
	quiet := flag.Bool("quiet", false, "quiet")

	flag.Parse()

	stage := env.Stage()

	var db *sql.DB
	var dialect string
	{
		env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
		if err != nil {
			log.Fatal(err)
		}
		conn, x, err := dbconnect.GetConnection(env)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		db = conn
		dialect = x
	}

	var app myapp.HTTPApp
	var appConfig myapp.Environment
	{
		env, err := myapp.GetEnvironment("app.yml", stage)
		if err != nil {
			log.Fatal(err)
		}
		appConfig = *env
		app, err = myapp.GetApp(env)
		if err != nil {
			log.Fatal(err)
		}
	}

	{
		app.Validator = validator.New()
		app.Validator.RegisterValidation("iffalse", func(f validator.FieldLevel) bool {
			rv := f.Parent().Elem().FieldByName(f.Param())
			if !rv.IsValid() {
				log.Fatalf("field %q not found", f.Param())
			}
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if rv.Interface().(bool) {
				return true
			}
			rs := f.Field()
			if rs.Kind() == reflect.Ptr {
				rs = rs.Elem()
			}
			s := rs.Interface().(string)
			return s != ""
		})
		app.Validator.RegisterStructValidation(
			myapp.CaptchaValidator(appConfig.CaptchaSolution),
			myapp.CaptchaInput{},
		)
	}

	if dialect == "sqlite3" || dialect == "mysql" {
		app.StepService = mysqlmodel.StepService{DB: db}
		app.ProtestService = mysqlmodel.ProtestService{DB: db}
		app.ContactMessageService = mysqlmodel.ContactMessageService{DB: db}
		// } else if dialect == "postgres" {
		// 	app.StepService = pgsqlmodel.StepService{DB: db}
		// 	app.ProtestService = pgsqlmodel.ProtestService{DB: db}
		// 	app.ContactMessageService = pgsqlmodel.ContactMessageService{DB: db}
	} else {
		panic(dialect)
	}

	r := mux.NewRouter()

	if appConfig.Host != "" {
		r = r.Host(appConfig.Host).Subrouter()
	}

	r.HandleFunc("/protests/by_author/{author_id}", app.GetProtests).Methods("GET")
	// r.HandleFunc("/protests/around", app.GetStepsAround).Methods("POST")
	r.HandleFunc("/protests/search", app.SearchProtests).Methods("POST")
	r.HandleFunc("/protests/create", app.CreateProtest).Methods("POST")
	r.HandleFunc("/protests/{id}", app.GetProtest).Methods("GET")
	r.HandleFunc("/protests/{id}", app.GetProtestWithPassword).Methods("POST")

	r.HandleFunc("/captcha/new", app.CaptchaNew).Methods("GET")
	r.Handle("/captcha/{id}.png", captcha.Server(captcha.StdWidth, captcha.StdHeight)).Methods("GET")

	r.HandleFunc("/contacts/create", app.CreateContactMessage).Methods("POST")

	verylimited := r.PathPrefix("/").Subrouter()
	verylimited.Use(myapp.GlobalRateLimiter(appConfig, 65536, 60, 0))
	r.HandleFunc("/geocode/search", myapp.Geocode(appConfig.GeoCoderCacheSize)).Methods("GET")

	veryprotected := r.PathPrefix("/").Subrouter()
	veryprotected.Use(myapp.RateLimiter(*appConfig.LoginRateLimit))
	veryprotected.HandleFunc("/admin/login", app.AdminLogin).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(myapp.Auth(appConfig))
	protected.HandleFunc("/contacts/delete/{id}", app.DeleteContactMessage).Methods("POST")
	protected.HandleFunc("/contacts/list", app.ListContactMessages).Methods("GET")

	{
		var fileSystem http.FileSystem = http.Dir("client/public")
		if appConfig.Statik {
			fileSystem = packr.NewBox("../client/public")
		}
		handler := http.FileServer(fileSystem)
		handler = myapp.InsertBefore(handler, "/", []byte("</body>"), myapp.CsrfToken)
		if !env.IsProd() {
			handler = myapp.InsertAfter(handler, "/", []byte("<head>"), myapp.JSErrorAlert)
		}
		handler = handlers.CompressHandler(handler)
		r.PathPrefix("/").Handler(handler)
	}

	canonicalhost := appConfig.CanonicalHost

	r.Use(myapp.Csrf(appConfig))
	var httpHandler http.Handler = r
	httpHandler = handlers.ContentTypeHandler(httpHandler, "application/json", "multipart/form-data")
	httpHandler = handlers.CORS()(httpHandler)
	httpHandler = handlers.CanonicalHost(canonicalhost, 302)(httpHandler)
	// httpHandler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(httpHandler)
	httpHandler = myapp.RateLimiter(*appConfig.GlobalRateLimit)(httpHandler)

	if *quiet == false {
		httpHandler = handlers.LoggingHandler(os.Stdout, httpHandler)
		// httpHandler = handlers.CombinedLoggingHandler(os.Stdout, httpHandler)
	}

	log.Println("Stage is ", stage)
	log.Println("Quiet=", *quiet)
	log.Printf("Serving %q using static assets %v\n", "client/public", appConfig.Statik)
	log.Printf("HTTP url: http://%v:%v\n", canonicalhost, appConfig.Port)
	log.Printf("Admin url: http://%v:%v/#!/admin?ia=true\n", canonicalhost, appConfig.Port)
	log.Printf("Admin key: %s\n", myapp.GetKey(appConfig))

	srv := &http.Server{
		Addr:    fmt.Sprint(":", appConfig.Port),
		Handler: httpHandler,
	}

	err := safeStart(srv.ListenAndServe)
	if err != nil {
		log.Fatal(err)
	}

	onSignal(os.Interrupt, func() {
		ctxSD, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		fmt.Printf("server shutdown err=%v\n", srv.Shutdown(ctxSD))
	})
}

func safeStart(h func() error) error {
	ferr := make(chan error)
	go func() {
		ferr <- h()
	}()
	select {
	case err := <-ferr:
		return err
	case <-time.After(time.Millisecond * 100):
	}
	return nil
}

func onSignal(s os.Signal, h func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, s)
	<-c
	h()
}
