package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const PORT = "8000"

var SharedApp App

func main() {
	var err error
	SharedApp.Router = Router{
		mux.NewRouter(),
	}
	SharedApp.HandleRoutes()	
	if err = SharedApp.InitDatabaseConnection(); err != nil {
		log.Println("UNABLE TO INIT DATABASE CONNECTION:", err)
		return
	}

	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
