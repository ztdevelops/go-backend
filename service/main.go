package main

import (
	"log"
	"net/http"
)

const PORT = "8000"

func main() {
	RouteHandlers()

	printString := "Listening for requests at http://localhost:" + PORT
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}