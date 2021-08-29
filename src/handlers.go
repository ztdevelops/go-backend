package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) HandleRoutes() {

	a.Router.HandleFunc("/", DefaultHandler).Methods((custom_types.RGET))
	a.Router.HandleFunc("/test", TestAPIHandler).Methods((custom_types.RGET))
	a.Router.HandleFunc("/signup", HandleSignUp).Methods((custom_types.RPOST))
	a.Router.HandleFunc("/signin", HandleSignIn).Methods((custom_types.RPOST))
	http.Handle("/", a.Router)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "default")
	fmt.Fprintf(w, "Default landing page.")
}

func TestAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "test")
	users := []custom_types.User{
		{Username: "user 1", Password: "123"},
		{Username: "user 2", Password: "456"},
	}
	json.NewEncoder(w).Encode(users)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented.")
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "sign up")
	received := &custom_types.User{}
	err := json.NewDecoder(r.Body).Decode(received); if err != nil {
		log.Println("Cannot decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(received.Password), 8)
	if err != nil {
		log.Println("Error hashing password:", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	submitted := &custom_types.User{
		Username: received.Username,
		Password: string(hashedPassword),
	}

	if err := SharedApp.SignUp(*submitted); err != nil {
		log.Println("Cannot sign up:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "sign in")
	received := &custom_types.User{}
	if err := json.NewDecoder(r.Body).Decode(received); err != nil {
		log.Println("Cannot decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := SharedApp.GetUser(received.Username)
	if err != nil {
		log.Println("Cannot get user:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(received.Password)); err != nil {
		log.Println("Comparison failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("Login OK.")
}