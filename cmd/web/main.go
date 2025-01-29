package main

import (
	"fmt"
	"net/http"
)

func main() {
	dep := initDependencies()
	//dep.DB.InsertReview("Горбатая гора (2005)", "вестерн, мелодрама, драма", "кино", "«Любовь - это сила природы»", 10, 1)
	defer dep.DB.CloseDB()
	mux := dep.routing()
	fmt.Println("Сервер запущен! Адрес: 127.0.0.1:8082")
	http.ListenAndServe("127.0.0.1:8082", mux)
}
