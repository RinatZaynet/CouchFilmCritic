package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
)

func (dep *Dependencies) reviewCreate(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.reviewCreate"
	const tmplCreateReview = "review-create.html"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodGet {
		logger.Warn("unsupported method. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	token, err := sesscookie.CheckCookie(r)
	if err != nil {
		logger.Warn("no session cookie. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if _, err := dep.JWT.CheckJWT(token); err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			logger.Warn("jwt-token expired. redirecting to login page")

			sesscookie.DeleteCookie(&w, r)

			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}

		logger.Error("failed to check jwt-token. redirecting to logout page", errslog.Err(err))

		http.Redirect(w, r, "/logout", http.StatusSeeOther)

		return
	}

	if err := dep.Templates.ExecuteTemplate(w, tmplCreateReview, nil); err != nil {
		logger.Error("failed to execute template", slog.String("tmpl", tmplCreateReview), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful of the handler work, execute template", slog.String("tmpl", tmplCreateReview))
}
