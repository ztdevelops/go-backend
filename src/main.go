package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var PORT = os.Getenv("PORT")

var SharedApp App

func main() {

	SharedApp.Router = Router{
		mux.NewRouter(),
	}
	SharedApp.HandleRoutes()
	SharedApp.InitDatabaseConnection()

	printString := "Listening for requests..."
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
