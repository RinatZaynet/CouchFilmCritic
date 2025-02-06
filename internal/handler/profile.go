package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/timefmt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	token, err := sesscookie.CheckCookie(r)
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sub, err := dep.JWT.CheckJWT(token)

	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) || errors.Is(err, auth.ErrInvalidToken) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	reviews, err := dep.DB.GetReviewsByAuthor(sub)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			dep.Templates.ExecuteTemplate(w, "profile.html", struct{ Reviews []*storage.Review }{reviews})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// локация должна соответствовать локации пользователя
	err = timefmt.TimeReviewsFmt(reviews, "Europe/Moscow")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	dep.Templates.ExecuteTemplate(w, "profile.html", struct{ Reviews []*storage.Review }{reviews})
}
