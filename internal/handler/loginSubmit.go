package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	hashpass "github.com/RinatZaynet/CouchFilmCritic/internal/hashingPassword"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/validation"
)

func (dep *Dependencies) loginSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.loginSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
		logger.Warn("unsupported method", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if _, err := sesscookie.CheckCookie(r); err == nil {
		logger.Warn("login submit attempt with existing session cookie")

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	nickname := r.FormValue("nickname")

	if !validation.IsValidNickname(nickname) {
		dep.Slogger.Info("not valid nickname", slog.String("nickname", nickname))

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	unique, err := dep.DB.IsUniqueNickname(nickname)
	if err != nil {
		logger.Error("failed to check nickname for uniqueness", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if unique {
		// переписать на алерт
		logger.Info("there is no user with such nickname", slog.String("nickname", nickname))

		return
	}

	if !validation.IsValidNickname(nickname) {
		dep.Slogger.Info("not valid nickname", slog.String("nickname", nickname))

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	pass := r.FormValue("password")

	if !validation.IsValidPassword(pass) {
		dep.Slogger.Info("not valid password", slog.String("nickname", nickname))

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	user, err := dep.DB.GetUserByNickname(nickname)
	if err != nil {
		logger.Error("failed to get user by nickname", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if err := dep.A2.CompareHashAndPassword([]byte(pass), user.PasswordHash); err != nil {
		if errors.Is(err, hashpass.ErrMismatchesTypes) {
			// переписать на алерт
			logger.Info("wrong password", slog.String("nickname", nickname))

			return
		}
		logger.Error("failed to hash and compare password", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	claims := &auth.Claims{
		Sub: nickname,
		Exp: time.Now().Add(240 * time.Hour).Unix(),
	}

	token, err := dep.JWT.GenJWT(claims)
	if err != nil {
		logger.Error("failed to generate jwt-token", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	sesscookie.CreateCookie(w, token)

	logger.Info("successful login user", slog.String("nickname", nickname))

	logger.Info("successful of the handler work")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
