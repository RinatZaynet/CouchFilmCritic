package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) reviewUpdate(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.reviewUpdate"
	const tmplReviewUpdate = "review-update.html"

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

	user, err := dep.DB.GetUserByNickname(nickname)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			logger.Warn("no user with this nickname",
				slog.String("nickname", nickname),
			)

			http.Redirect(w, r, "/profile", http.StatusSeeOther)

			return
		}

		logger.Error("failed to get user by nickname",
			slog.String("nickname", nickname),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		logger.Warn("the variable id was not found in the request",
			slog.String("nickname", nickname),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	reviewID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to conv string to int",
			slog.String("nickname", user.Nickname),
			slog.String("val", id),
			errslog.Err(err),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	review, err := dep.DB.GetReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			logger.Warn("no review with this id",
				slog.String("nickname", nickname),
				slog.Int("reviewID", reviewID),
			)

			http.Redirect(w, r, "/profile", http.StatusSeeOther)

			return
		}

		logger.Error("failed to get review by id",
			slog.Int("id", reviewID),
			slog.String("nickname", nickname),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if review.Author != user.Nickname {
		logger.Warn("attempt to update a review by a user who did not write it",
			slog.String("nickname", nickname),
			slog.Int("reviewID", reviewID),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	if err = dep.Templates.ExecuteTemplate(w, tmplReviewUpdate, struct{ Review *storage.Review }{review}); err != nil {
		logger.Error("failed to execute template", slog.String("tmpl", tmplReviewUpdate), errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful of the handler work, execute template", slog.String("tmpl", tmplReviewUpdate))
}
