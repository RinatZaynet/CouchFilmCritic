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

func (dep *Dependencies) reviewCreateSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.reviewCreateSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
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

	ratingStr := r.FormValue("rating")

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		logger.Error("failed to conv string to int. redirecting to profile page",
			slog.String("nickname", user.NickName),
			slog.String("val", ratingStr),
			errslog.Err(err),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	// добавить валидацию с алертами
	id, err := dep.DB.InsertReview(
		r.FormValue("work_title"),
		r.FormValue("genres"),
		r.FormValue("work_type"),
		r.FormValue("review"),
		rating,
		user.ID,
	)

	if err != nil {
		logger.Error("failed to insert new review",
			slog.String("nickname", user.NickName),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful create new review",
		slog.String("nickname", user.NickName),
		slog.Int("reviewID", id),
	)

	logger.Info("successful of the handler work, redirecting to profile page")

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
