package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/timefmt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) profile(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.profile"
	const tmplProfile = "profile.html"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodGet {
		logger.Warn("unsupported method", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	token, err := sesscookie.CheckCookie(r)
	if err != nil {
		logger.Warn("no session cookie", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	sub, err := dep.JWT.CheckJWT(token)
	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			logger.Info("jwt-token expired")

			sesscookie.DeleteCookie(w, r)

			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}

		logger.Error("failed to check jwt-token", errslog.Err(err))

		http.Redirect(w, r, "/logout", http.StatusSeeOther)

		return
	}

	reviews, err := dep.DB.GetReviewsByAuthor(sub)
	if err != nil {
		logger.Error("failed to get reviews", slog.String("nickname", sub), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}
	// локация должна соответствовать локации пользователя
	if err := timefmt.TimeReviewsFmt(reviews, "Europe/Moscow"); err != nil {
		logger.Error("failed to format time", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if err := dep.Templates.ExecuteTemplate(w, tmplProfile, struct{ Reviews []*storage.Review }{reviews}); err != nil {
		logger.Error("failed to execute template", slog.String("tmpl", tmplProfile), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful of the handler work, execute template", slog.String("tmpl", tmplProfile))
}
