package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ztdevelops/go-project/src/helpers/custom_structs"
)

const ENDPOINT_HIT = "ENDPOINT HIT:"

func RouteHandlers() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/test", testAPIHandler)

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(ENDPOINT_HIT, "default")
	fmt.Fprintf(w, "Default landing page.")
}

func testAPIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(ENDPOINT_HIT, "test")
	users := []custom_structs.User{
		{Username: "user 1", Password: "123"},
		{Username: "user 2", Password: "456"},
	}
	json.NewEncoder(w).Encode(users)
}