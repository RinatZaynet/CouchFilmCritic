package handler

import (
	"net/http"
	"strconv"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

func (dep *Dependencies) reviewCreateSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	token, err := sesscookie.CheckCookie(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sub, err := dep.JWT.CheckJWT(token)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	user, err := dep.DB.GetUserByNickName(sub)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// добавить валидацию с алертами
	_, err = dep.DB.InsertReview(
		r.FormValue("work_title"), r.FormValue("genres"), r.FormValue("work_type"), r.FormValue("review"), rating, user.ID)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
