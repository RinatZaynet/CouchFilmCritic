package main

import (
	"log"
	"net/http"
)

func (dep *dependencies) regPage(w http.ResponseWriter, r *http.Request) {
	err := dep.Templates.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (dep *dependencies) mainPage(w http.ResponseWriter, r *http.Request) {
	err := dep.Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
