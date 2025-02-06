package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

func (dep *Dependencies) createReview(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := dep.Templates.Execute(w, "create-review.html")
		if err != nil {
			log.Fatal(err)
		}
	}
	if r.Method == http.MethodPost {
		rating, err := strconv.Atoi(r.FormValue("rating"))
		if err != nil {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}

		token, err := sesscookie.CheckCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		sub, err := dep.JWT.CheckJWT(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		user, err := dep.DB.GetUserByNickName(sub)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		// добавить валидацию с алертами
		dep.DB.InsertReview(r.FormValue("work_title"), r.FormValue("genres"), r.FormValue("work_type"), r.FormValue("review"), rating, user.ID)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
