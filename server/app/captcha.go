package app

import (
	"encoding/json"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gorilla/csrf"
	validator "gopkg.in/go-playground/validator.v9"
)

type CaptchaInput struct {
	CaptchaID       string `json:"captchaid" validate:"required" conform:"alphanum"`
	CaptchaSolution string `json:"captchasolution" validate:"required" conform:"alphanum"`
}

func (c CaptchaInput) GetCaptchaID() string       { return c.CaptchaID }
func (c CaptchaInput) GetCaptchaSolution() string { return c.CaptchaSolution }

type captchaer interface {
	GetCaptchaID() string
	GetCaptchaSolution() string
}

func CaptchaValidator(magicSolution string) func(sl validator.StructLevel) {
	return func(sl validator.StructLevel) {
		d := sl.Current().Interface().(captchaer)
		if magicSolution != "" && d.GetCaptchaSolution() == magicSolution {
			return
		}
		if !captcha.VerifyString(d.GetCaptchaID(), d.GetCaptchaSolution()) {
			sl.ReportError(d.GetCaptchaSolution(), "captchaSolution", "captchaSolution", "captcha", "")
		}
	}
}

//CaptchaNew starts a new captcha challenge.
func (h HTTPApp) CaptchaNew(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	json.NewEncoder(w).Encode(map[string]string{
		"id": captcha.New(),
	})
}
