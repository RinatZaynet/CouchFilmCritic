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
)

func (dep *Dependencies) loginSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.loginSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
		logger.Warn("unsupported method. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if _, err := sesscookie.CheckCookie(r); err == nil {
		logger.Warn("login submit attempt with existing session cookie. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	nickName := r.FormValue("nickname")

	unique, err := dep.DB.IsNickNameUnique(nickName)
	if err != nil {
		logger.Error("failed to check nickname for uniqueness", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if unique {
		// переписать на алерт
		logger.Warn("there is no user with such nickname", slog.String("nickname", nickName))

		return
	}

	user, err := dep.DB.GetUserByNickName(nickName)
	if err != nil {
		logger.Error("failed to get user by nickname", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if err := dep.A2.CompareHashAndPassword([]byte(r.FormValue("password")), user.PasswordHash); err != nil {
		if errors.Is(err, hashpass.ErrMismatchesTypes) {
			// переписать на алерт
			logger.Warn("wrong password", slog.String("nickname", nickName))

			return
		}
		logger.Error("failed to hash and compare password", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	claims := &auth.Claims{
		Sub: nickName,
		Exp: time.Now().Add(240 * time.Hour).Unix(),
	}

	token, err := dep.JWT.GenJWT(claims)
	if err != nil {
		logger.Error("failed to generate jwt-token", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	sesscookie.CreateCookie(&w, token)

	logger.Info("successful login user", slog.String("nickname", nickName))

	logger.Info("successful of the handler work, redirecting to index page")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
