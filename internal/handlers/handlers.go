package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
)

func (dep *Dependencies) index(w http.ResponseWriter, r *http.Request) {
	/*reviews, err := dep.DB.GetReviewsByAuthor("Rinat")
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
	err = dep.Templates.ExecuteTemplate(w, "main.html", struct{ Reviews []*storage.Review }{reviews})
	if err != nil {
		fmt.Println(err)
	}*/
	/*
		reviews, err := dep.DB.GetLatestReviews()
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		// локация должна соответствовать локации пользователя
		err = formatTimeReviews(reviews, "Europe/Moscow")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		err = dep.Templates.ExecuteTemplate(w, "main.html", struct{ Reviews []*models.Review }{reviews})
		if err != nil {
			fmt.Println(err)
		}*/

	/*err := dep.Templates.ExecuteTemplate(w, "main.html", nil)
	if err != nil {
		log.Fatal(err)
	}*/
	token, err := sesscookie.CheckCookie(r)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	sub, err := dep.JWT.CheckJWT(token)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	user, err := dep.DB.GetUserByNickName(sub)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintln(w, user.ID, user.NickName, user.Email, user.PasswordHash, user.SignUpDate.In(loc))

	/*_, err := dep.DB.InsertUser("Rinat", "rinat@mail.ru", "13r1jgfu9cxcvx6vspmz")
	if err != nil {
		log.Fatal(err)
	}
	user, err := dep.DB.GetUserByNickName("Rinat")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, user)
	//dep.deleteSessCookie(r, &w)*/
}

func (dep *Dependencies) login(w http.ResponseWriter, r *http.Request) {
	err := dep.Templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (dep *Dependencies) reg(w http.ResponseWriter, r *http.Request) {
	err := dep.Templates.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (dep *Dependencies) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Http method %s is incorrect, use method %s. Status: %d.", r.Method, http.MethodPost, http.StatusMethodNotAllowed)
		return
	}
	// TODO: Написать функцию валидатор входящих от пользователя данных и перенести туда весь код ниже

	nickName := r.FormValue("nickname")
	unique, err := dep.DB.IsNickNameUnique(nickName)
	if err != nil {
		//http.Redirect(w, r, "/reg", http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	if !unique {
		// Переписать чтобы высвечивался алерт о том что пользователь с таким nickname уже существует
		http.Redirect(w, r, "/reg", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	unique, err = dep.DB.IsEmailUnique(email)
	if err != nil {
		//http.Redirect(w, r, "/reg", http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	if !unique {
		// Переписать чтобы высвечивался алерт о том что пользователь с таким email уже существует
		http.Redirect(w, r, "/reg", http.StatusSeeOther)
		return
	}
	if r.FormValue("password0") != r.FormValue("password1") {
		http.Redirect(w, r, "/reg", http.StatusSeeOther)
		return
	}
	_, err = dep.DB.InsertUser(nickName, email, r.FormValue("password0"))
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	claims := &auth.Claims{
		Sub: nickName,
		Exp: time.Now().Add(240 * time.Hour).Unix(),
	}

	token, err := dep.JWT.GenJWT(claims)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	sesscookie.CreateCookie(&w, token)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
