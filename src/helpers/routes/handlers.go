package routes

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
)

func (r *Router) HandleRoutes() {

	r.HandleFunc("/", DefaultHandler).Methods((RGET))
	r.HandleFunc("/test", TestAPIHandler).Methods((RGET))
	http.Handle("/", r)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(ENDPOINT_HIT, "default")
	fmt.Fprintf(w, "Default landing page.")
}

func TestAPIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(ENDPOINT_HIT, "test")
	users := []custom_types.User{
		{Username: "user 1", Password: "123"},
		{Username: "user 2", Password: "456"},
	}
	json.NewEncoder(w).Encode(users)
}