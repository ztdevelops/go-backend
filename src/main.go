package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ztdevelops/go-project/src/helpers/database"
	"github.com/ztdevelops/go-project/src/helpers/routes"
)

const PORT = "8000"

type App struct {
	Router 	routes.Router
	DB		database.Database
}

func main() {
	app := App{}
	app.Router = routes.Router{
		mux.NewRouter(),
	}
	app.Router.HandleRoutes()	
	app.DB.InitDatabaseConnection()

	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
