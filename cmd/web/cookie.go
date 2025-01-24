package main

import (
	"net/http"
	"time"
)

func (dep *dependencies) checkSessCookie(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	if cookie.Expires.Before(time.Now()) {
		return false
	}
	return true
}

func (dep *dependencies) createSessCookie(w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   JWTtoken,
		Expires: time.Now().Add(240 * time.Hour),
	}
	http.SetCookie(w, cookie)
}
