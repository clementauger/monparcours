package app

import (
	"bytes"
	"log"
	"net/http"
	"os"

	appconf "github.com/clementauger/monparcours/server/config/app"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
	"github.com/tomasen/realip"
)

//Csrf middleware
func Csrf(config appconf.Environment) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(config.CsrfKey),
		csrf.RequestHeader("X-tomate"),
		csrf.CookieName("ctomate"),
		csrf.FieldName("ftomate"),
		csrf.Secure(false),
	)
}

//CsrfToken of a request to string
func CsrfToken(r *http.Request) string {
	return string(csrf.TemplateField(r))
}

//JSErrorAlert template for a given request.
func JSErrorAlert(r *http.Request) string {
	return `<script>
	  window.addEventListener("error", handleError, true);
	  function handleError(evt) {
	      if (evt.message) {
	        alert("error: "+evt.message +" at linenumber: "+evt.lineno+" of file: "+evt.filename);
	      } else {
	        alert("error: "+evt.type+" from element: "+(evt.srcElement || evt.target));
	      }
	  }
	  </script>`
}

type writeReplacer struct {
	http.ResponseWriter
	search  []byte
	replace func(*http.Request) string
	buf     []byte
	r       *http.Request
	dir     string
}

func (w *writeReplacer) Write(in []byte) (int, error) {
	if w.buf == nil {
		w.buf = []byte{}
	}
	w.buf = append(w.buf, in...)
	n := len(in)
	if index := bytes.LastIndex(w.buf, w.search); index > -1 {
		var r []byte
		if w.dir == "before" {
			g := []byte(w.replace(w.r))
			n += len(g)
			r = append(g, w.buf[index:]...)
			w.buf = append(w.buf[:index], r...)
		} else {
			g := []byte(w.replace(w.r))
			n += len(g)
			r = append(r, w.buf[:index+len(w.search)]...)
			r = append(r, g...)
			r = append(r, w.buf[index:]...)
			w.buf = r
		}
	}
	return n, nil
}

func (w *writeReplacer) Flush() {
	w.ResponseWriter.Write(w.buf[:])
	w.buf = w.buf[:0]
}

func (w *writeReplacer) WriteHeader(statusCode int) {
	w.Header().Del("Content-length")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.ResponseWriter.WriteHeader(statusCode)
}

//InsertAfter middleware to inject html after some token.
func InsertAfter(h http.Handler, path string, search []byte, replace func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			replacer := &writeReplacer{ResponseWriter: w, search: search, replace: replace, r: r, dir: "after"}
			defer func() {
				replacer.Flush()
			}()
			h.ServeHTTP(replacer, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

//InsertBefore middleware to inject html before some token.
func InsertBefore(h http.Handler, path string, search []byte, replace func(*http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			replacer := &writeReplacer{ResponseWriter: w, search: search, replace: replace, r: r, dir: "before"}
			defer func() {
				replacer.Flush()
			}()
			h.ServeHTTP(replacer, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

//RateLimiter middleware that varies by path and remoteaddr.
func RateLimiter(config appconf.RateLimit) func(http.Handler) http.Handler {
	store, err := memstore.New(config.Size)
	if err != nil {
		log.Fatal(err)
	}

	quota := throttled.RateQuota{MaxRate: throttled.PerMin(config.RPM), MaxBurst: config.Burst}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true, RemoteAddr: true},
	}
	return httpRateLimiter.RateLimit
}

//GlobalRateLimiter middleware that does not very by.
func GlobalRateLimiter(config appconf.Environment, size, rpm, burst int) func(http.Handler) http.Handler {
	store, err := memstore.New(size)
	if err != nil {
		log.Fatal(err)
	}

	quota := throttled.RateQuota{MaxRate: throttled.PerMin(rpm), MaxBurst: burst}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{},
	}
	return httpRateLimiter.RateLimit
}

//HTPPLog middleware to log http requests.
func HTPPLog() func(http.Handler) http.Handler {
	return func(httpHandler http.Handler) http.Handler {
		return realIP{Next: handlers.CombinedLoggingHandler(os.Stdout, httpHandler)}
	}
}

type realIP struct {
	Next http.Handler
}

func (rip realIP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	r.RemoteAddr = realip.FromRequest(r)
	rip.Next.ServeHTTP(w, r)
	r.RemoteAddr = ip
}
