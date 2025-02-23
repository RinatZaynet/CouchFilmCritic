package auth

import (
	"net/http"
	"time"
)

func GetAuthToken(r *http.Request) (authToken string) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func CreateAuthCookie(w http.ResponseWriter, authToken string) {
	cookie := &http.Cookie{
		Name:    "auth_token",
		Value:   authToken,
		Path:    "/",
		Expires: time.Now().Add(240 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func DeleteAuthCookie(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth_token")
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name: "auth_token",
		Path: "/",
		//MaxAge: -1,
		Expires: time.Now().Add(-1),
	}

	http.SetCookie(w, cookie)
}
