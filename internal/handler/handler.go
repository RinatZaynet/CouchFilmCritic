package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth"
	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/sesscookie"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/timefmt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

func (dep *Dependencies) index(w http.ResponseWriter, r *http.Request) {
	reviews, err := dep.DB.GetReviewsByAuthor("Rinat")
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
		err = dep.Templates.ExecuteTemplate(w, "index.html", struct{ Reviews []*models.Review }{reviews})
		if err != nil {
			fmt.Println(err)
		}*/

	/*err := dep.Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}*/
	/*
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
		fmt.Fprintln(w, user.ID, user.NickName, user.Email, user.PasswordHash, user.SignUpDate.In(loc)) */

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
	if r.Method == http.MethodGet {
		err := dep.Templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Fatal(err)
		}

		return
	}
	if r.Method == http.MethodPost {
		unique, err := dep.DB.IsNickNameUnique(r.FormValue("nickname"))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		if !unique {
			// переписать на алерт
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		// Добавить хеширование пароля и проверку на соответствия пользователя с таким ником и паролем в БД
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

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

func (dep *Dependencies) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (dep *Dependencies) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	token, err := sesscookie.CheckCookie(r)
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sub, err := dep.JWT.CheckJWT(token)

	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) || errors.Is(err, auth.ErrInvalidToken) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	reviews, err := dep.DB.GetReviewsByAuthor(sub)

	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			dep.Templates.ExecuteTemplate(w, "profile.html", struct{ Reviews []*storage.Review }{reviews})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dep.Templates.ExecuteTemplate(w, "profile.html", struct{ Reviews []*storage.Review }{reviews})
}
