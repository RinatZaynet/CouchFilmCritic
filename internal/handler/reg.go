package handler

import (
	"log"
	"net/http"
)

func (dep *Dependencies) reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := dep.Templates.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
