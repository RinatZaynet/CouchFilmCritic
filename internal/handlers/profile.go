package handlers

import (
	"log/slog"
	"net/http"

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

	sub, err := dep.checkAuth(w, r)

	if err != nil {
		logger.Error("failed to check auth", errslog.Err(err))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if sub == "" {
		logger.Info("no session cookie")

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	nickname := sub

	reviews, err := dep.DB.GetReviewsByAuthor(nickname)

	if err != nil {
		logger.Error("failed to get reviews", slog.String("nickname", nickname), errslog.Err(err))

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
