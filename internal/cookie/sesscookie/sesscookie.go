package sesscookie

import (
	"fmt"
	"net/http"
	"time"
)

func CheckCookie(r *http.Request) (token string, err error) {
	const fn = "cookie.sesscookie.CheckCookie"

	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return cookie.Value, nil
}

func CreateCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:    "jwt_token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(240 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("jwt_token")
	if err != nil {
		return
	}
	cookie := &http.Cookie{
		Name:   "jwt_token",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
