package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func RegPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("./ui/html/*"))
	err := tmpl.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("./ui/html/*"))
	w.Header().Add("X-Content-Type-Options", "nosniff")
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/reg", RegPage)
	mux.HandleFunc("/main", MainPage)

	fmt.Println("Сервер запущен! Адрес: 127.0.0.1:8082")
	http.ListenAndServe("127.0.0.1:8082", mux)
}
