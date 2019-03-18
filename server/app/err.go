package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/clementauger/monparcours/server/env"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

//SqlError discriminates errors kind.
type SqlError interface {
	Sql() string
}

//UserError discriminates errors kind.
type UserError interface {
	User()
}

//AuthError for a generic user auth error report.
type AuthError struct {
	error
}

//User discriminates errors kind.
func (v AuthError) User() {}

//InputError for a generic user input error report.
type InputError struct {
	error
}

//User discriminates errors kind.
func (v InputError) User() {}

//ValidationError json encode the validator errors.
type ValidationError struct {
	errs validator.ValidationErrors
}

//User discriminates errors kind.
func (v ValidationError) User() {}

func (v ValidationError) Error() string {
	out := map[string]string{}
	for _, field := range v.errs {
		key := strings.ToLower(field.StructNamespace())
		rule := field.ActualTag()
		out[key] = rule
	}
	x, _ := json.Marshal(out)
	return fmt.Sprintf("%s", x)
}

func HandleValidationError(w http.ResponseWriter, err error) error {
	if x, ok := err.(validator.ValidationErrors); ok {
		http.Error(w, ValidationError{errs: x}.Error(), http.StatusBadRequest)
		return nil

	}
	return err
}
func HandleUserError(w http.ResponseWriter, err error) error {
	if _, ok := err.(UserError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil

	}
	return err
}
func HandleOtherErrors(w http.ResponseWriter, err error) error {
	httpMsg := err.Error()
	if env.IsProd() {
		httpMsg = "unrecoverable error"
	}
	longErr := httpMsg
	switch x := errors.Cause(err).(type) {
	case SqlError:
		longErr = fmt.Sprintf("%v\n%v", x.Sql(), httpMsg)
	default:
	}
	if err, ok := err.(stackTracer); ok {
		longErr += fmt.Sprintln()
		for _, f := range err.StackTrace() {
			longErr += fmt.Sprintf("%+s:%d\n", f, f)
		}
	}
	log.Println(longErr)
	http.Error(w, httpMsg, http.StatusInternalServerError)
	return nil
}
