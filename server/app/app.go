package app

import (
	"bytes"
	"log"
	"net/http"
	"regexp"

	"github.com/clementauger/monparcours/server/model"
	"github.com/gorilla/csrf"
	"github.com/leebenson/conform"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
	validator "gopkg.in/go-playground/validator.v9"
)

//HTTPApp hosts routes implementation.
type HTTPApp struct {
	Env                   Environment
	ProtestService        model.ProtestService
	StepService           model.StepService
	ContactMessageService model.ContactMessageService
	Validator             *validator.Validate
}

func init() {
	conform.AddSanitizer("text", text)
	conform.AddSanitizer("alphanum", alphanum)
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

func Csrf(config Environment) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(config.CsrfKey),
		csrf.RequestHeader("X-tomate"),
		csrf.CookieName("ctomate"),
		csrf.FieldName("ftomate"),
		csrf.Secure(false),
	)
}
func CsrfToken(r *http.Request) string {
	return string(csrf.TemplateField(r))
}
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

func RateLimiter(config RateLimit) func(http.Handler) http.Handler {
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

func GlobalRateLimiter(config Environment, size, rpm, burst int) func(http.Handler) http.Handler {
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
