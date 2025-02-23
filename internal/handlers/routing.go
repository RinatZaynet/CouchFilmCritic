package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/hashpass/argon2"
	"github.com/RinatZaynet/CouchFilmCritic/internal/jwtutill"
	mwLogger "github.com/RinatZaynet/CouchFilmCritic/internal/middleware/logger"
	"github.com/RinatZaynet/CouchFilmCritic/internal/middleware/requestid"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage/mysql"
)

type Dependencies struct {
	Templates *template.Template
	DB        *mysql.ManagerDB
	JWT       *jwtutill.Manager
	A2        *argon2.Manager
	Slogger   *slog.Logger
}

func (dep *Dependencies) Routing() *http.ServeMux {
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

	loggerMux := http.NewServeMux()
	mwL := mwLogger.New(dep.Slogger)
	loggerMux.Handle("/", mwL(mux))

	reqIDMux := http.NewServeMux()
	mwR := requestid.New(dep.Slogger, dep.JWT)
	reqIDMux.Handle("/", mwR(loggerMux))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return reqIDMux
}
