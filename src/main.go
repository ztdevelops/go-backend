package main

import (
	"log"
	"net/http"

	"github.com/ztdevelops/go-project/src/helpers/database" 
)

const PORT = "8000"

func main() {
	RouteHandlers()
	database.InitDatabaseConnection()

	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}