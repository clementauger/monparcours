package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	myapp "github.com/clementauger/monparcours/server/app"
	"github.com/clementauger/monparcours/server/env"
	"github.com/clementauger/monparcours/server/service"
	"github.com/gobuffalo/packr"
	"github.com/gocraft/health"

	"github.com/dchest/captcha"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Save the stream as a global variable
var stream = health.NewStream()

//ServeHTTP application.
func ServeHTTP(ctx context.Context) {
	quiet := flag.Bool("quiet", false, "quiet")

	flag.Parse()

	{
		sink := health.NewJsonPollingSink(time.Minute, time.Minute*20)
		stream.AddSink(sink)
		sink.StartServer("localhost:5020")
	}

	stage := env.Stage()

	srvs := service.Service{}
	if err := srvs.Init(stage); err != nil {
		log.Fatal(err)
	}
	defer srvs.Close()

	var app myapp.HTTPApp
	if err := app.Init(stage, &srvs); err != nil {
		log.Fatal(err)
	}
	appConfig := app.Env

	r := &myapp.Router{
		Stream:       stream,
		Router:       mux.NewRouter(),
		ErrorHandler: myapp.HandleHTTPError,
	}

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

	r.HandleFunc("/rgpd", app.RGPD).Methods("GET")

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
	if canonicalhost != "" {
		httpHandler = handlers.CanonicalHost(canonicalhost, 302)(httpHandler)
	}
	httpHandler = myapp.RateLimiter(*appConfig.GlobalRateLimit)(httpHandler)

	if *quiet == false {
		httpHandler = myapp.HTPPLog()(httpHandler)
	}

	log.Println("Stage is ", stage)
	log.Println("Quiet=", *quiet)
	log.Printf("Serving %q using static assets %v\n", "client/public", appConfig.Statik)
	log.Printf("HTTP url: http://%v:%v\n", canonicalhost, appConfig.Port)
	log.Printf("Admin url: http://%v:%v/#!/admin?ia=true\n", canonicalhost, appConfig.Port)
	log.Printf("Admin key: %s\n", myapp.GetKey(appConfig))

	srv := &http.Server{
		Addr:         fmt.Sprint(":", appConfig.Port),
		Handler:      httpHandler,
		ReadTimeout:  *appConfig.ReadTimeout,
		WriteTimeout: *appConfig.WriteTimeut,
	}

	err := safeStart(srv.ListenAndServe)
	if err != nil {
		log.Fatal("failed to start http server", err)
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
