package handler

import (
	"fmt"
	"net/http"
)

func (dep *Dependencies) userSubmit(w http.ResponseWriter, r *http.Request) {
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

	hashedPass, err := dep.A2.HashingPassword([]byte(r.FormValue("password0")))
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	_, err = dep.DB.InsertUser(nickName, email, hashedPass)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
