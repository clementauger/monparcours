package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/clementauger/monparcours/server/env"
	"gopkg.in/go-playground/validator.v9"
)

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

//HandleHTTPError writes err content on w
// if it is an UserError, it responds with a StatusBadRequest code,
// if err is non nil, and the environement is not production, it responds with an StatusInternalServerError code,
// otherwise it writes a dummy message with a response code StatusInternalServerError,
// and writes the error to stderr.
func HandleHTTPError(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}
	if x, ok := err.(validator.ValidationErrors); ok {
		http.Error(w, ValidationError{errs: x}.Error(), http.StatusBadRequest)
		return nil

	} else if _, ok := err.(UserError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil

	} else if !env.IsProd() {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	http.Error(w, "unrecoverable error", http.StatusInternalServerError)
	return err
}
