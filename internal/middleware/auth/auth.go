package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

var noAuth = map[string]struct{}{
	"/":             {},
	"/reg":          {},
	"/reg/submit":   {},
	"/login":        {},
	"/login/submit": {},
	"/logout":       {},
}

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/auth"),
		)

		log.Info("auth middleware enabled")

		handler := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := noAuth[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			token, err := sesscookie.CheckCookie(r)

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			if tParts := strings.Split(token, "."); len(tParts) != 3 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}
