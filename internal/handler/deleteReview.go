package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) deleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Println("123")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
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
	fmt.Println("1233")
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
