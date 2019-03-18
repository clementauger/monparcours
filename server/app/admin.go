package app

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/clementauger/httpr"
	appconf "github.com/clementauger/monparcours/server/config/app"
	"github.com/clementauger/st"
	sth "github.com/clementauger/st/http"
	"github.com/gorilla/csrf"
)

type authHandler struct {
	next http.Handler
	Env  appconf.Environment
}

func (a authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	iak, err := r.Cookie("iak")
	if err != nil {
		return AuthError{fmt.Errorf("failed to login")}
	}
	if !checkToken(a.Env, iak.Value, 0) && !checkToken(a.Env, iak.Value, time.Hour*24*-1) {
		return AuthError{fmt.Errorf("failed to login")}
	}
	a.next.ServeHTTP(w, r)
	return nil
}

//Auth middleware to require logged user
func Auth(env appconf.Environment) func(http.Handler) httpr.ErrHandler {
	return func(h http.Handler) httpr.ErrHandler {
		return authHandler{next: h, Env: env}
	}
}

func getToken(env appconf.Environment, d time.Duration) string {
	k := fmt.Sprintf("%v:%v:%v",
		env.AdminKey, env.AdminSalt,
		time.Now().Add(d).Format("2006-01-02"),
	)
	k = fmt.Sprintf("%x", md5.Sum([]byte(k)))
	return k
}

func checkToken(env appconf.Environment, token string, d time.Duration) bool {
	return token == getToken(env, d)
}

//GetKey for admin access.
func GetKey(env appconf.Environment) string {
	k := fmt.Sprintf("%v:%v",
		env.AdminKey, env.AdminSalt,
	)
	k = fmt.Sprintf("%x", md5.Sum([]byte(k)))
	return k
}

func checkKey(env appconf.Environment, token string) bool {
	return token == GetKey(env)
}

type adminLoginInput struct {
	Key string `json:"key" validate:"required,max=60" conform:"alphanum"`
}

//AdminLogin attempts to resolve admin login challenge.
func (h HTTPApp) AdminLogin(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input adminLoginInput

	return st.
		Map(sth.Decode(&input, sth.JSONDecode(r))).
		Map(sth.Conform(input)).
		Map(sth.Validate(input, h.Validator)).
		Map(func(i adminLoginInput) (bool, error) {
			if !checkKey(h.Env, i.Key) {
				return false, AuthError{fmt.Errorf("failed to login")}
			}
			cookie := &http.Cookie{
				Name:     "iak",
				Value:    getToken(h.Env, 0),
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour * 6),
				Path:     "/",
			}
			http.SetCookie(w, cookie)
			return true, nil
		}).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}
