package handler

import (
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/validation"
)

func (dep *Dependencies) regSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.regSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
		logger.Warn("unsupported method", slog.String("method", r.Method))

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	nickname := r.FormValue("nickname")

	if !validation.IsValidNickname(nickname) {
		dep.Slogger.Info("not valid nickname", slog.String("nickname", nickname))

		http.Redirect(w, r, "/reg", http.StatusSeeOther)

		return
	}

	unique, err := dep.DB.IsUniqueNickname(nickname)
	if err != nil {
		logger.Error("failed to check nickname for uniqueness", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if !unique {
		// переписать на алерт
		logger.Warn("this nickname is already taken", slog.String("nickname", nickname))

		return
	}

	email := r.FormValue("email")

	if !validation.IsValidEmail(email) {
		dep.Slogger.Info("not valid email", slog.String("nickname", nickname))

		http.Redirect(w, r, "/reg", http.StatusSeeOther)

		return
	}

	unique, err = dep.DB.IsUniqueEmail(email)
	if err != nil {
		logger.Error("failed to check email for uniqueness", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if !unique {
		// переписать на алерт
		logger.Warn("this email is already taken", slog.String("nickname", nickname))

		return
	}

	pass := r.FormValue("password")
	passConfirm := r.FormValue("passwordConfirm")

	if !validation.IsValidNewPassword(pass, passConfirm) {
		dep.Slogger.Info("not valid password", slog.String("nickname", nickname))

		http.Redirect(w, r, "/reg", http.StatusSeeOther)

		return
	}

	hashedPass := dep.A2.HashingPassword([]byte(r.FormValue("password")))

	_, err = dep.DB.InsertUser(nickname, email, hashedPass)
	if err != nil {
		logger.Error("failed to insert user", errslog.Err(err))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	logger.Info("successful registration user", slog.String("nickname", nickname))

	logger.Info("successful of the handler work")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
