package handler

import (
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

func (dep *Dependencies) logout(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.logout"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
		logger.Warn("unsupported method", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	sesscookie.DeleteCookie(w, r)

	logger.Info("successful logout user")

	logger.Info("successful of the handler work")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
