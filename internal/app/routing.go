package app

import "net/http"

func (dep *dependencies) routing() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dep.helpers.mainPage)
	mux.HandleFunc("/login", dep.loginPage)
	mux.HandleFunc("/reg", dep.regPage)
	mux.HandleFunc("/create/user", dep.createUser)
	return mux
}
