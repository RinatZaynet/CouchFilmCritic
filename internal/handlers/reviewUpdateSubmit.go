package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) reviewUpdateSubmit(w http.ResponseWriter, r *http.Request) {
	const fn = "handler.reviewUpdateSubmit"

	logger := dep.Slogger.With(slog.String("func", fn))

	logger.Info("start of the handler work")

	if r.Method != http.MethodPost {
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
	_ = user
}
