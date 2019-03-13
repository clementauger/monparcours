package app

import (
	"net/http"

	"github.com/clementauger/monparcours/server/model"
	"github.com/clementauger/st"
	sth "github.com/clementauger/st/http"
	"github.com/gorilla/csrf"
)

type contactMessageInput struct {
	CaptchaInput `validate:"required,dive,required"`
	model.ContactMessage
}

//CreateContactMessage decodes the body as a json request, validates the input data,
// writes the database, then respond the written object.
func (h HTTPApp) CreateContactMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input contactMessageInput

	return st.
		Map(sth.Decode(&input, sth.JSONDecode(r))).
		Map(sth.Conform(input)).
		Map(sth.Validate(input, h.Validator)).
		Map(func(i contactMessageInput) (model.ContactMessage, error) {
			return h.Services.ContactMessage.Insert(i.ContactMessage)
		}).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

type listContactMessagesInput struct {
	Offset int64 `schema:"offset"`
	Limit  int64 `schema:"limit"`
}

//ListContactMessages ...
func (h HTTPApp) ListContactMessages(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input listContactMessagesInput

	return st.
		Map(sth.Decode(&input, sth.GetDecode(r))).
		Map(func(input listContactMessagesInput) ([]model.ContactMessage, error) {
			return h.Services.ContactMessage.GetAll()
		}).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

type deleteContactMessagesInput struct {
	ID int64 `schema:"id"`
}

//DeleteContactMessage ...
func (h HTTPApp) DeleteContactMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input deleteContactMessagesInput

	return st.
		Map(sth.Decode(&input, sth.MuxDecode(r))).
		Map(func(input deleteContactMessagesInput) (bool, error) {
			d, err := h.Services.ContactMessage.Get(input.ID)
			if err == nil {
				err = h.Services.ContactMessage.Delete(d)
			}
			return err == nil, err
		}).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}
