package handler

import (
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/timefmt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) index(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.index"
	const tmplIndex = "index.html"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodGet {
		logger.Warn("unsupported method", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	reviews, err := dep.DB.GetLatestReviews()
	if err != nil {
		logger.Error("failed to get latest reviews", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	// локация должна соответствовать локации пользователя
	if err := timefmt.TimeReviewsFmt(reviews, "Europe/Moscow"); err != nil {
		logger.Error("failed to format time", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if err := dep.Templates.ExecuteTemplate(w, tmplIndex, struct{ Reviews []*storage.Review }{reviews}); err != nil {
		logger.Error("failed to execute template", slog.String("tmpl", tmplIndex), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful of the handler work, execute template", slog.String("tmpl", tmplIndex))
}
