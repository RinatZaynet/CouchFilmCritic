package handler

import (
	"html/template"
	"log/slog"
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
	Slogger   *slog.Logger
}

func Routing(dep *Dependencies) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dep.index)
	mux.HandleFunc("/login", dep.login)
	mux.HandleFunc("/login/submit", dep.loginSubmit)
	mux.HandleFunc("/reg", dep.reg)
	mux.HandleFunc("/reg/submit", dep.regSubmit)
	mux.HandleFunc("/profile", dep.profile)
	mux.HandleFunc("/logout", dep.logout)
	mux.HandleFunc("/review/create", dep.reviewCreate)
	mux.HandleFunc("/review/create/submit", dep.reviewCreateSubmit)
	mux.HandleFunc("/review/update/", dep.reviewUpdate)
	mux.HandleFunc("/review/delete/", dep.reviewDelete)

	authMux := http.NewServeMux()
	authMux.Handle("/", authmiddleware.AuthMid(mux))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	authMux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return authMux
}
