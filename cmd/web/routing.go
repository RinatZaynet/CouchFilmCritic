package main

import "net/http"

func (dep *dependencies) routing() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", dep.mainPage)
	mux.HandleFunc("/reg", dep.regPage)
	mux.HandleFunc("/login", dep.loginPage)
	return mux
}
