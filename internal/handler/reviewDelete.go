package handler

import (
	"net/http"
	"strconv"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) reviewDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	token, err := sesscookie.CheckCookie(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = dep.JWT.CheckJWT(token)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reviewID, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	err = dep.DB.DeleteReviewByID(reviewID)
	if err != nil {
		if err == storage.ErrNoRows {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
