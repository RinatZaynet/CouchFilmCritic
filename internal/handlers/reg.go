package handlers

import (
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
)

func (dep *Dependencies) reg(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.reg"
	const tmplReg = "registration.html"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodGet {
		logger.Warn("unsupported method", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	sub, err := dep.checkAuth(w, r)

	if err != nil {
		logger.Error("failed to check auth", errslog.Err(err))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if sub != "" {
		nickname := sub

		logger.Warn("registration attempt with existing session cookie", slog.String("nickname", nickname))
	}

	if err := dep.Templates.ExecuteTemplate(w, tmplReg, nil); err != nil {
		logger.Error("failed to execute template", slog.String("tmpl", tmplReg), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful of the handler work, execute template", slog.String("tmpl", tmplReg))
}
