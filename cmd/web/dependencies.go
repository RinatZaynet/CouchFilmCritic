package main

import (
	"log"
	"text/template"

	"github.com/RinatZaynet/CouchFilmCritic/pkg/models/jwt"
	"github.com/RinatZaynet/CouchFilmCritic/pkg/models/mysql"
)

type dependencies struct {
	Templates *template.Template
	DB        *mysql.ManagerDB
	JWT       *jwt.ManagerJWT
}

// Внедрить конфигурацию
func initDependencies() *dependencies {
	tmpl := ParseTemplates("./ui/html/*")

	db, err := mysql.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	clientJWT, err := jwt.NewClientJWT()
	if err != nil {
		log.Fatal(err)
	}

	dep := &dependencies{
		Templates: tmpl,
		DB:        &mysql.ManagerDB{Database: db},
		JWT:       clientJWT,
	}
	return dep
}
