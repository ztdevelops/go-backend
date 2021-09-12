package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func main() {

	SharedApp.Router = Router{
		mux.NewRouter(),
	}
	SharedApp.HandleRoutes()	
	SharedApp.InitDatabaseConnection()

	allowedHeaders := []string{"GET", "POST"}
	allowedMethods := []string{"Content-Type", "Origin", "Accept", "*"}
	corsWrapper := GetCorsWrapper(allowedHeaders, allowedMethods)
	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":"+PORT, corsWrapper.Handler(SharedApp.Router)))
}
