package app

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	sth "github.com/clementauger/st/http"

	"github.com/clementauger/monparcours/server/model"

	"github.com/clementauger/st"
	"github.com/gorilla/csrf"
)

//CreateProtest decodes the body as a json request, validates the input data,
// writes the database, then respond the written object.
func (h HTTPApp) CreateProtest(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input model.Protest
	return st.
		Map(sth.Decode(&input, sth.JSONDecode(r))).
		Map(sth.Conform(input)).
		Map(sth.Validate(input, h.Validator)).
		Map(func(m model.Protest) model.Protest {
			if m.Password != "" {
				hasher := md5.New()
				hasher.Write([]byte(m.Password + ":" + h.Env.PwdSalt))
				m.Password = fmt.Sprintf("%x", hasher.Sum(nil))
			}
			return m
		}).
		Map(h.Services.Protest.Insert).
		Map(h.Services.Step.InsertSteps).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

type getProtestInput struct {
	ID int64 `schema:"id"`
}

//GetProtest by its ID.
func (h HTTPApp) GetProtest(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input getProtestInput

	return st.
		Map(sth.Decode(&input, sth.MuxDecode(r))).
		Map(func(input getProtestInput) (model.Protest, error) {
			return h.Services.Protest.Get(input.ID)
		}).
		Map(h.Services.Step.GetSteps).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

type getProtestPwdInput struct {
	ID       int64  `schema:"id"`
	Password string `schema:"pwd" json:"pwd"`
}

//GetProtestWithPassword by its ID and password.
func (h HTTPApp) GetProtestWithPassword(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input getProtestPwdInput

	return st.
		Map(sth.Decode(&input, sth.MuxDecode(r), sth.JSONDecode(r))).
		Map(func(i getProtestPwdInput) getProtestPwdInput {
			hasher := md5.New()
			hasher.Write([]byte(i.Password + ":" + h.Env.PwdSalt))
			i.Password = fmt.Sprintf("%x", hasher.Sum(nil))
			return i
		}).
		Map(func(input getProtestPwdInput) (model.Protest, error) {
			return h.Services.Protest.GetWithPassword(input.ID, input.Password)
		}).
		Map(h.Services.Step.GetProtectedSteps).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

type getProtestsInput struct {
	AuthorID string `schema:"author_id"`
}

//GetProtests by their authorID.
func (h HTTPApp) GetProtests(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input getProtestsInput

	return st.
		Map(sth.Decode(&input, sth.MuxDecode(r))).
		Map(func(input getProtestsInput) ([]model.Protest, error) {
			return h.Services.Protest.GetByAuthorID(input.AuthorID)
		}).
		Map(st.Each(h.Services.Step.GetSteps)).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

type searchProtestInput struct {
	Title     *string    `schema:"title" json:"title"`
	Protest   *string    `schema:"protest" json:"protest"`
	Organizer *string    `schema:"organizer" json:"organizer"`
	AtLat     *float64   `schema:"lat" json:"lat"`
	AtLng     *float64   `schema:"lng" json:"lng"`
	StartDate *time.Time `schema:"date_start" json:"date_start"`
	EndDate   *time.Time `schema:"date_end" json:"date_end"`
}

//SearchProtests a location and a date.
func (h HTTPApp) SearchProtests(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	var input searchProtestInput

	return st.
		Map(sth.Decode(&input, sth.JSONDecode(r))).
		Map(func(input searchProtestInput) ([]model.Protest, error) {
			return h.Services.Protest.SearchProtests(
				input.Title, input.Protest, input.Organizer,
				input.StartDate, input.EndDate,
				input.AtLat, input.AtLng,
				50.0,
			)
		}).
		Map(st.Each(h.Services.Step.GetSteps)).
		Map(sth.JSONEncode(w)).
		Sink()

	// HandleHTTPError(w, err)
}

//ProtestInterest updates the protest interest count.
func (h HTTPApp) ProtestInterest(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	// var err error
	// HandleHTTPError(w, err)
	return nil
}

//ProtestView updates the protest view count.
func (h HTTPApp) ProtestView(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

	// var err error
	// HandleHTTPError(w, err)
	return nil
}
