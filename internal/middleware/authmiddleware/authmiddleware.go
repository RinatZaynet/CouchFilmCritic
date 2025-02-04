package authmiddleware

import (
	"net/http"
	"strings"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

var noAuth = map[string]struct{}{
	"/":      {},
	"/reg":   {},
	"/login": {},
}

func AuthMid(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuth[r.URL.Path]; ok {
			h.ServeHTTP(w, r)
			return
		}
		token, err := sesscookie.CheckCookie(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		if tParts := strings.Split(token, "."); len(tParts) != 3 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		h.ServeHTTP(w, r)
	})
}
