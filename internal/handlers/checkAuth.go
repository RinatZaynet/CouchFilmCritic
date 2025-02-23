package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/jwtutill"
)

func (dep *Dependencies) checkAuth(w http.ResponseWriter, r *http.Request) (sub string, err error) {
	const fn = "handlers.dependencies.checkAuth"

	logger := dep.Slogger.With(slog.String("func", fn))

	token := auth.GetAuthToken(r)
	if token == "" {
		auth.DeleteAuthCookie(w, r)

		return "", nil
	}

	sub, err = dep.JWT.Parse(token)

	if err != nil {
		if errors.Is(err, jwtutill.ErrTokenExpired) {
			logger.Info("jwt-token expired")

			auth.DeleteAuthCookie(w, r)

			return "", nil
		}

		return "", err
	}

	return sub, nil
}
