package app

import (
	"log"
	"reflect"
	"regexp"

	appconf "github.com/clementauger/monparcours/server/config/app"
	"github.com/clementauger/monparcours/server/service"
	"github.com/leebenson/conform"
	validator "gopkg.in/go-playground/validator.v9"
)

//HTTPApp hosts routes implementation.
type HTTPApp struct {
	Env       appconf.Environment
	Services  *service.Service
	Validator *validator.Validate
}

//Init http application
func (h *HTTPApp) Init(stage string, srvs *service.Service) error {

	appConfig, err := appconf.GetEnvironment("app.yml", stage)
	if err != nil {
		log.Fatal(err)
	}
	h.Env = *appConfig

	{
		h.Validator = validator.New()
		h.Validator.RegisterValidation("iffalse", iffalse)
		h.Validator.RegisterStructValidation(
			CaptchaValidator(appConfig.CaptchaSolution),
			CaptchaInput{},
		)
	}

	{
		conform.AddSanitizer("text", text)
		conform.AddSanitizer("alphanum", alphanum)
	}

	h.Services = srvs

	return nil
}

func iffalse(f validator.FieldLevel) bool {
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
}

var patterns = map[string]*regexp.Regexp{
	"alphanum": regexp.MustCompile("[^0-9\\pL]"),
	"text":     regexp.MustCompile("[^0-9\\pL!?.,;:()[]\"'+-=*%]"),
}

func text(s string) string {
	return patterns["text"].ReplaceAllLiteralString(s, "")
}
func alphanum(s string) string {
	return patterns["alphanum"].ReplaceAllLiteralString(s, "")
}
