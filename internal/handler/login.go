package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	hashpass "github.com/RinatZaynet/CouchFilmCritic/internal/hashingPassword"
)

func (dep *Dependencies) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := dep.Templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Fatal(err)
		}

		return
	}
	if r.Method == http.MethodPost {
		nickName := r.FormValue("nickname")

		unique, err := dep.DB.IsNickNameUnique(nickName)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		if unique {
			// переписать на алерт
			fmt.Fprintf(w, "%s", err)
			return
		}
		user, err := dep.DB.GetUserByNickName(nickName)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		err = dep.A2.CompareHashAndPassword([]byte(r.FormValue("password")), user.PasswordHash)
		if err != nil {
			if errors.Is(err, hashpass.ErrMismatchesTypes) {
				// переписать на алерт
				fmt.Fprintf(w, "%s", err)
				return
			}
			fmt.Fprintf(w, "%s", err)
			return
		}

		claims := &auth.Claims{
			Sub: nickName,
			Exp: time.Now().Add(240 * time.Hour).Unix()}
		token, err := dep.JWT.GenJWT(claims)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		sesscookie.CreateCookie(&w, token)

	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
