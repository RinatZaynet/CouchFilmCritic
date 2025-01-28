package main

import (
	"net/http"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/pkg/models"
)

func (dep *dependencies) checkSessCookie(r *http.Request) (sub string, err error) {
	sub = ""
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		return sub, err
	}
	sub, err = dep.JWT.CheckTokenJWT(cookie.Value)
	if cookie.Expires.Before(time.Now()) {
		return sub, err
	}
	return sub, nil
}

func (dep *dependencies) createSessCookie(w *http.ResponseWriter, sub string) (err error) {
	token, err := dep.JWT.GenTokenJWT(&models.Session{Sub: sub})
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    "jwt_token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(240 * time.Hour),
	}
	http.SetCookie(*w, cookie)
	return nil
}

func (dep *dependencies) deleteSessCookie(r *http.Request, w *http.ResponseWriter) {
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		return
	}
	if cookie.Expires.Before(time.Now()) {
		return
	}
	cookie.Expires = time.Now().AddDate(0, 0, 1)

	http.SetCookie(*w, cookie)
}
