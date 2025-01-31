package main

import (
	"github.com/RinatZaynet/CouchFilmCritic/internal/app"
)

func main() {
	/*cfg := config.MustConfigParsing()
	dep := initDependencies()
	defer dep.DB.CloseDB()
	mux := dep.routing()
	fmt.Println("Сервер запущен! Адрес: 127.0.0.1:8081")
	http.ListenAndServe("127.0.0.1:8081", mux)*/
	app.Run()
}
