package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) reviewUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("1")
	token, err := sesscookie.CheckCookie(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("2")
	sub, err := dep.JWT.CheckJWT(token)
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}
	fmt.Println("3")
	user, err := dep.DB.GetUserByNickName(sub)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("4")
	review, err := dep.DB.GetReviewByID(user.ID)
	fmt.Println("5")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("6")
	err = dep.Templates.ExecuteTemplate(w, "review-update.html", struct{ Review *storage.Review }{review})
	fmt.Println("7")
	// Переписать
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("8")
	return
}
