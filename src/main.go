package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const PORT = "8000"

var SharedApp App

func main() {

	SharedApp.Router = Router{
		mux.NewRouter(),
	}
	SharedApp.HandleRoutes()	
	SharedApp.InitDatabaseConnection()

	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
