package handler

import (
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
)

func (dep *Dependencies) regSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.regSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
		logger.Warn("unsupported method. redirecting to index page", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}
	// TODO: Написать функцию валидатор входящих от пользователя данных и перенести туда весь код ниже

	nickName := r.FormValue("nickname")

	unique, err := dep.DB.IsNickNameUnique(nickName)
	if err != nil {
		logger.Error("failed to check nickname for uniqueness", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if unique {
		// переписать на алерт
		logger.Warn("this nickname is already taken", slog.String("nickname", nickName))

		return
	}

	email := r.FormValue("email")
	unique, err = dep.DB.IsEmailUnique(email)
	if err != nil {
		logger.Error("failed to check email for uniqueness", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if unique {
		// переписать на алерт
		logger.Warn("this email is already taken", slog.String("nickname", nickName))

		return
	}

	if r.FormValue("password0") != r.FormValue("password1") {
		// переписать на алерт
		logger.Warn("mismatch passwords", slog.String("nickname", nickName))

		return
	}

	hashedPass, err := dep.A2.HashingPassword([]byte(r.FormValue("password0")))
	if err != nil {
		logger.Error("failed to hashing password", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	_, err = dep.DB.InsertUser(nickName, email, hashedPass)
	if err != nil {
		logger.Error("failed to insert user", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful registration user", slog.String("nickname", nickName))

	logger.Info("successful of the handler work, redirecting to login page")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
