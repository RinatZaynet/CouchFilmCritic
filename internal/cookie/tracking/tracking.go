package tracking

import (
	"net/http"
	"time"
)

func GetRequestID(r *http.Request) string {
	cookie, err := r.Cookie("request_id")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func CreateRequestIDCookie(w http.ResponseWriter, requestID string) {
	cookie := &http.Cookie{
		Name:    "request_id",
		Value:   requestID,
		Path:    "/",
		Expires: time.Now().Add(240 * time.Hour),
	}

	http.SetCookie(w, cookie)
}
