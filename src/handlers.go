package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
)

// HandleRoutes initialises the connections to all the explicitly coded routes.
func (a *App) HandleRoutes() {
	app, err := InitFirebase()
	if err != nil {
		log.Fatal("failed to init firebase:", err)
	}

	a.App = *app
	a.Router.HandleFunc("/", DefaultHandler).Methods((custom_types.RGET))
	a.Router.HandleFunc("/test", TestAPIHandler).Methods((custom_types.RGET))

	// APIs (No need for auth)
	a.Router.HandleFunc("/api/signup", NotImplemented).Methods((custom_types.RPOST))
	a.Router.HandleFunc("/api/signin", LogInWithFirebase).Methods((custom_types.RPOST))

	// APIs (Require auth)
	a.Router.HandleFunc("/api/test", a.TestVerifyToken).Methods(custom_types.RGET)
	http.Handle("/", a.Router)
}

func (a *App) TestVerifyToken(w http.ResponseWriter, r *http.Request) {
	if err := VerifyToken(&a.App, r); err != nil {
		return
	}
	log.Println("token ok.")
}

// Respond responds to the API requests with a status code and message.
func (w *Writer) Respond(code int, message interface{}) {
	r := Response{
		Status:  code,
		Message: message,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(r)
}

// SetContentType sets the content type of the response.
func (w *Writer) SetContentType(ct string) {
	w.Header().Set("Content-Type", ct)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "default")
	fmt.Fprintf(w, "Default landing page.")
}

func TestAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "test")
	users := []custom_types.User{
	}
	json.NewEncoder(w).Encode(users)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented.")
}
