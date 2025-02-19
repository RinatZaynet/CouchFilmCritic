package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/validation"
)

func (dep *Dependencies) reviewCreateSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.reviewCreateSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
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

	user, err := dep.DB.GetUserByNickname(sub)
	if err != nil {
		logger.Error("failed to get user by nickname",
			slog.String("nickname", sub),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	workTitle := r.FormValue("work_title")

	if !validation.IsValidWorkTitle(workTitle) {
		dep.Slogger.Info("not valid work title",
			slog.String("nickname", sub),
			slog.String("work title", workTitle))

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	genres := r.Form["genres"]

	fmtGenres := strings.Join(genres, ", ")

	if !validation.IsValidGenres(genres) {
		dep.Slogger.Info("not valid genres",
			slog.String("nickname", sub),
			slog.String("genres", fmtGenres))

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	workType := r.FormValue("work_type")

	if !validation.IsValidWorkType(workType) {
		dep.Slogger.Info("not valid work type",
			slog.String("nickname", sub),
			slog.String("work type", workType))

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	review := r.FormValue("review")

	if !validation.IsValidReview(review) {
		dep.Slogger.Info("not valid review",
			slog.String("nickname", sub),
			slog.String("review", review))

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	ratingStr := r.FormValue("rating")

	if !validation.IsValidRating(ratingStr) {
		dep.Slogger.Info("not valid rating",
			slog.String("nickname", sub),
			slog.String("rating", ratingStr))

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		logger.Error("failed to conv string to int",
			slog.String("nickname", user.Nickname),
			slog.String("val", ratingStr),
			errslog.Err(err),
		)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

		return
	}

	id, err := dep.DB.InsertReview(
		workTitle,
		fmtGenres,
		workType,
		review,
		rating,
		user.Nickname,
	)

	if err != nil {
		logger.Error("failed to insert new review",
			slog.String("nickname", user.Nickname),
			errslog.Err(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful create new review",
		slog.String("nickname", user.Nickname),
		slog.Int("reviewID", id),
	)

	logger.Info("successful of the handler work")

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
