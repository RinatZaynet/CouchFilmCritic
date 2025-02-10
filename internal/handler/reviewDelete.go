package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
)

func (dep *Dependencies) reviewDelete(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.reviewDelete"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodDelete {
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

	sub, err := dep.JWT.CheckJWT(token)
	if err != nil {
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

	user, err := dep.DB.GetUserByNickName(sub)
	if err != nil {
		logger.Error("failed to get user by nickname",
			slog.String("nickname", sub),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		logger.Warn("the variable id was not found in the request. redirecting to profile page",
			slog.String("nickname", sub),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	reviewID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to conv string to int. redirecting to profile page",
			slog.String("nickname", user.NickName),
			slog.String("val", id),
			errslog.Err(err),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	review, err := dep.DB.GetReviewByID(reviewID)
	if err != nil {
		logger.Error("failed to get review by id",
			slog.Int("reviewID", reviewID),
			slog.String("nickname", sub),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if review.Author != user.NickName {
		logger.Warn("attempt to delete a review by a user who did not write it. redirecting to profile page",
			slog.String("nickname", sub),
			slog.Int("reviewID", reviewID),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	if err := dep.DB.DeleteReviewByID(reviewID); err != nil {
		logger.Error("failed to delete review by id",
			slog.String("nickname", user.NickName),
			slog.String("id", id),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful delete review",
		slog.String("nickname", user.NickName),
		slog.Int("reviewID", reviewID),
	)

	logger.Info("successful of the handler work, redirecting to profile page")

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
