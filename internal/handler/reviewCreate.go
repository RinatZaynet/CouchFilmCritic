package handler

import (
	"log"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

func (dep *Dependencies) reviewCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	err = dep.Templates.ExecuteTemplate(w, "review-create.html", nil)
	// Переписать
	if err != nil {
		log.Fatal(err)
	}
}
