package handler

import (
	"html/template"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth/jwt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/hashingPassword/argon2"
	"github.com/RinatZaynet/CouchFilmCritic/internal/middleware/authmiddleware"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage/mysql"
)

type Dependencies struct {
	Templates *template.Template
	DB        *mysql.ManagerDB
	JWT       *jwt.ManagerJWT
	A2        *argon2.ManagerArgon2
}

func Routing(dep *Dependencies) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dep.index)
	mux.HandleFunc("/login", dep.login)
	mux.HandleFunc("/reg", dep.reg)
	mux.HandleFunc("/create/user", dep.createUser)
	mux.HandleFunc("/profile", dep.profile)
	mux.HandleFunc("/logout", dep.logout)
	mux.HandleFunc("/create/review", dep.createReview)
	mux.HandleFunc("/delete/review/", dep.deleteReview)

	authMux := http.NewServeMux()
	authMux.Handle("/", authmiddleware.AuthMid(mux))
	return authMux
}
