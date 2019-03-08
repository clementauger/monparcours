package app

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/clementauger/st"
	sth "github.com/clementauger/st/http"
	"github.com/gorilla/csrf"
)

type authHandler struct {
	http.Handler
	Env Environment
}

func (a authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	iak, err := r.Cookie("iak")
	if err != nil {
		HandleHTTPError(w, AuthError{fmt.Errorf("failed to login")})
		return
	}
	if !checkToken(a.Env, iak.Value, 0) && !checkToken(a.Env, iak.Value, time.Hour*24*-1) {
		HandleHTTPError(w, AuthError{fmt.Errorf("failed to login")})
		return
	}
	a.Handler.ServeHTTP(w, r)
}

func Auth(env Environment) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return authHandler{Handler: h, Env: env}
	}
}

func getToken(env Environment, d time.Duration) string {
	k := fmt.Sprintf("%v:%v:%v",
		env.AdminKey, env.AdminSalt,
		time.Now().Add(d).Format("2006-01-02"),
	)
	k = fmt.Sprintf("%x", md5.Sum([]byte(k)))
	return k
}

func checkToken(env Environment, token string, d time.Duration) bool {
	return token == getToken(env, d)
}

//GetKey for admin access.
func GetKey(env Environment) string {
	k := fmt.Sprintf("%v:%v",
		env.AdminKey, env.AdminSalt,
	)
	k = fmt.Sprintf("%x", md5.Sum([]byte(k)))
	return k
}

func checkKey(env Environment, token string) bool {
	return token == GetKey(env)
}

type adminLoginInput struct {
	Key string `json:"key" validate:"required,max=60" conform:"alphanum"`
}

//AdminLogin attempts to resolve admin login challenge.
func (h HTTPApp) AdminLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input adminLoginInput

	err := st.
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

	HandleHTTPError(w, err)
}
