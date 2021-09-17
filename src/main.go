package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ztdevelops/go-project/src/lib/custom"
	"github.com/ztdevelops/go-project/src/lib/middleware"
)


func main() {
	if err := custom.LoadEnv(); err != nil {
		log.Println("error loading .env:", err)
	}

	SharedApp.Router = Router{
		mux.NewRouter(),
	}
	SharedApp.HandleRoutes()
	SharedApp.InitDatabaseConnection()

	allowedHeaders := []string{"GET", "POST"}
	allowedMethods := []string{"Content-Type", "Origin", "Accept", "*"}
	corsWrapper := middleware.GetCorsWrapper(allowedHeaders, allowedMethods)
	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":"+PORT, corsWrapper.Handler(SharedApp.Router)))
}
