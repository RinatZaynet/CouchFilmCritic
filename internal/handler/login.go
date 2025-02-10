package handler

import (
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
)

func (dep *Dependencies) login(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.login"
	const tmplLogin = "login.html"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodGet {
		logger.Warn("unsupported method. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if _, err := sesscookie.CheckCookie(r); err == nil {
		logger.Warn("login attempt with existing session cookie. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if err := dep.Templates.ExecuteTemplate(w, tmplLogin, nil); err != nil {
		logger.Error("failed to execute template", slog.String("tmpl", tmplLogin), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful of the handler work, execute template", slog.String("tmpl", tmplLogin))
}
