package app

import (
	"net/http"
	"time"
)

//RGPD deletes cookies of an user visiting this handler.
func (h HTTPApp) RGPD(w http.ResponseWriter, r *http.Request) error {
	for _, c := range r.Cookies() {
		c.Expires = time.Unix(0, 0)
		http.SetCookie(w, c)
	}
	w.Header().Set("Location", "https://giphy.com/explore/be-gone")
	w.WriteHeader(http.StatusMovedPermanently)

	return nil
}
