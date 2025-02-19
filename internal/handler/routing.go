package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth/jwt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/hashingPassword/argon2"
	mwAuth "github.com/RinatZaynet/CouchFilmCritic/internal/middleware/auth"
	mwLogger "github.com/RinatZaynet/CouchFilmCritic/internal/middleware/logger"
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

	newMWAuth := mwAuth.New(dep.Slogger)
	authMux.Handle("/", newMWAuth(mux))

	loggerMux := http.NewServeMux()

	newMWLogger := mwLogger.New(dep.Slogger)
	loggerMux.Handle("/", newMWLogger(authMux))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	loggerMux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return loggerMux
}
