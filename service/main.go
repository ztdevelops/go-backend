package main

import (
	"io"
	"log"
	"net/http"
)

const PORT = "8000"

func main() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello World!... \n")
	}

	http.HandleFunc("/hello", handler)
	printString := "Listening for requests at http://localhost:" + PORT + "/hello" 
	log.Println(printString)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
}