package main

import (
	"log"
	"net/http"
)

func (dep *dependencies) mainPage(w http.ResponseWriter, r *http.Request) {
	err := dep.Templates.ExecuteTemplate(w, "index.html", nil)
	//_, err := dep.DB.InsertUser("Rinat", "rinat@mail.ru", "13r1jgfu9cxcvx6vspmz")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//user, err := dep.DB.GetUserByNickName("Rinat")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Fprintln(w, user)
}
func (dep *dependencies) regPage(w http.ResponseWriter, r *http.Request) {
	err := dep.Templates.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
func (dep *dependencies) loginPage(w http.ResponseWriter, r *http.Request) {

}
