package app

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	lru "github.com/hashicorp/golang-lru"
)

func Geocode(size int) func(w http.ResponseWriter, r *http.Request) error {
	l, err := lru.New(size)
	if err != nil {
		log.Fatal(err)
	}
	url, _ := url.Parse("https://nominatim.openstreetmap.org/")

	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(w http.ResponseWriter, r *http.Request) error {

		rr, _ := http.NewRequest("GET", r.URL.String(), nil)
		rr.URL.Path = "/search"
		rr.URL.Host = url.Host
		rr.URL.Scheme = url.Scheme
		rr.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		rr.Header.Set("User-Agent", r.Header.Get("User-Agent"))
		rr.Header.Set("Referer", r.Header.Get("Referer"))
		rr.Header.Set("Accept", r.Header.Get("Accept"))
		rr.Header.Set("Accept-Encoding", r.Header.Get("Accept-Encoding"))
		rr.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))
		rr.Host = url.Host

		key := rr.URL.String()
		v, ok := l.Get(key)
		if ok {
			w.Write(v.([]byte))
			return nil
		}
		x := &resp{ResponseWriter: w}
		proxy.ServeHTTP(x, rr)
		l.Add(key, x.buf)
		return nil
	}
}

type resp struct {
	http.ResponseWriter
	buf []byte
}

func (r *resp) Write(d []byte) (int, error) {
	r.buf = append(r.buf, d...)
	return r.ResponseWriter.Write(d)
}
