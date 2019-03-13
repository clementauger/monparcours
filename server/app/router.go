package app

import (
	"net/http"
	"time"

	"github.com/gocraft/health"
	"github.com/gorilla/mux"
)

//Router handles handler errors
type Router struct {
	*mux.Router
	Stream       *health.Stream
	ErrorHandler func(http.ResponseWriter, error) error
}

//HandleFunc handles a router func
func (re *Router) HandleFunc(path string, of func(http.ResponseWriter, *http.Request) error) *Route {
	f := func(w http.ResponseWriter, r *http.Request) {
		job := re.Stream.NewJob(path)
		startTime := time.Now()

		var originalErr error
		var finalErr error

		originalErr = of(w, r)
		defer job.EventErr(path, originalErr)
		if re.ErrorHandler != nil {
			finalErr = re.ErrorHandler(w, originalErr)
		}
		job.Timing(path, time.Since(startTime).Nanoseconds())

		if finalErr == nil {
			if originalErr == nil {
				job.Complete(health.Success)
			} else {
				job.Complete(health.ValidationError)
			}
		} else {
			job.Complete(health.Error)
		}
	}
	return &Route{Route: re.Router.HandleFunc(path, f), ErrorHandler: re.ErrorHandler, Stream: re.Stream}
}

//Handle an http handler
func (re *Router) Handle(path string, handler http.Handler) *Route {
	return &Route{Route: re.Router.Handle(path, handler), ErrorHandler: re.ErrorHandler, Stream: re.Stream}
}

//PathPrefix a route
func (re *Router) PathPrefix(u string) *Route {
	return &Route{Route: re.Router.PathPrefix(u), ErrorHandler: re.ErrorHandler, Stream: re.Stream}
}

//Host of the router
func (re *Router) Host(u string) *Route {
	return &Route{Route: re.Router.Host(u), ErrorHandler: re.ErrorHandler, Stream: re.Stream}
}

//Route of the router
type Route struct {
	*mux.Route
	Stream       *health.Stream
	ErrorHandler func(http.ResponseWriter, error) error
}

//Subrouter to handle handlers with error
func (re *Route) Subrouter() *Router {
	return &Router{Router: re.Route.Subrouter(), ErrorHandler: re.ErrorHandler, Stream: re.Stream}
}
