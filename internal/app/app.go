package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RinatZaynet/CouchFilmCritic/internal/auth/jwt"
	"github.com/RinatZaynet/CouchFilmCritic/internal/config"
	"github.com/RinatZaynet/CouchFilmCritic/internal/handlers"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/errslog"
	"github.com/RinatZaynet/CouchFilmCritic/internal/storage/mysql"
)

func Run() {
	cfg := config.MustConfigParsing()
	log := initLogger(cfg.Env)

	log.Debug("debug messages are enabled")
	log.Info("info messages are enabled")

	log.Info("logger successful initialization", slog.String("env", cfg.Env))

	tmpl, err := parseTemplates(cfg.TemplatesPath)
	if err != nil {
		log.Error("failed to parse templates", errslog.Err(err))
		os.Exit(1)
	}
	db, err := mysql.OpenDB(cfg.Dsn)
	if err != nil {
		log.Error("failed to connect to mysql database", errslog.Err(err))
		os.Exit(1)
	}
	defer db.Close()
	clientJWT, err := jwt.NewClientJWT(cfg.JWTSecret)

	if err != nil {
		log.Error("failed to init jwt", errslog.Err(err))
		os.Exit(1)
	}

	dep := &handlers.Dependencies{
		Templates: tmpl,
		DB:        &mysql.ManagerDB{Database: db},
		JWT:       clientJWT,
	}

	mux := handlers.Routing(dep)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	// доделать
	log.Info("server start", slog.String("addr", cfg.Address))
	err = server.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", errslog.Err(err))
	}
}
