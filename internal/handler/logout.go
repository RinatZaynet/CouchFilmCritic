package handler

import (
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

func (dep *Dependencies) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sesscookie.DeleteCookie(r, &w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
