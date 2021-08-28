package main

import (
	"log"
	"net/http"

	"github.com/ztdevelops/go-project/src/helpers" 
)

const PORT = "8000"

func main() {
	RouteHandlers()
	helpers.InitDatabaseConnection()

	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}