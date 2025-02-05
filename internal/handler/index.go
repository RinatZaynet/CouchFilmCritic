package handler

import (
	"fmt"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/timefmt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) index(w http.ResponseWriter, r *http.Request) {
	reviews, err := dep.DB.GetLatestReviews()
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	// локация должна соответствовать локации пользователя
	err = timefmt.TimeReviewsFmt(reviews, "Europe/Moscow")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	err = dep.Templates.ExecuteTemplate(w, "index.html", struct{ Reviews []*storage.Review }{reviews})
	if err != nil {
		fmt.Println(err)
	}
}
