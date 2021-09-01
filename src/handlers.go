package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
	"golang.org/x/crypto/bcrypt"
)

// HandleRoutes initialises the connections to all the explicitly coded routes.
func (a *App) HandleRoutes() {
	a.Router.HandleFunc("/", DefaultHandler).Methods((custom_types.RGET))
	a.Router.HandleFunc("/test", TestAPIHandler).Methods((custom_types.RGET))

	// APIs
	a.Router.HandleFunc("/api/signup", HandleSignUp).Methods((custom_types.RPOST))
	a.Router.HandleFunc("/api/signin", HandleSignIn).Methods((custom_types.RPOST))
	a.Router.HandleFunc("/api/posts", NotImplemented).Methods(custom_types.RGET, custom_types.RPUT, custom_types.RPOST, custom_types.RDELETE)
	a.Router.HandleFunc("/api/posts/:id", NotImplemented).Methods(custom_types.RGET, custom_types.RPUT, custom_types.RPOST, custom_types.RDELETE)
	http.Handle("/", a.Router)
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
		{Username: "user 1", Password: "123"},
		{Username: "user 2", Password: "456"},
	}
	json.NewEncoder(w).Encode(users)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented.")
}

// HandleSignUp creates an account based on the submitted request.
func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "sign up")
	writer := Writer{w}
	writer.SetContentType(ContentTypeJSON)

	received := &custom_types.User{}
	err := json.NewDecoder(r.Body).Decode(received)
	if err != nil {
		log.Println("Cannot decode:", err)
		writer.Respond(http.StatusBadRequest, err.Error())
		return
	}

	if received.Username == "" || received.Password == "" {
		log.Println("Received empty credentials. Rejecting request.")
		writer.Respond(http.StatusBadRequest, "Failed to create account: invalid credentials")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(received.Password), 8)
	if err != nil {
		log.Println("Error hashing password:", err)
		writer.Respond(http.StatusBadRequest, err.Error())
		return
	}

	submitted := &custom_types.User{
		Username: received.Username,
		Password: string(hashedPassword),
	}

	if err := SharedApp.SignUp(*submitted); err != nil {
		log.Println("Cannot sign up:", err)
		writer.Respond(http.StatusBadRequest, err.Error())
		return
	}

	writer.Respond(http.StatusOK, "User successfully created.")
}

// HandleSignIn attempts to log the user into their accounts based on the submitted credentials.
func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "sign in")
	writer := Writer{w}
	writer.SetContentType(ContentTypeJSON)

	received := &custom_types.User{}
	if err := json.NewDecoder(r.Body).Decode(received); err != nil {
		writer.Respond(http.StatusBadRequest, err.Error())
		return
	}
	user, err := SharedApp.GetUser(received.Username)
	if err != nil {
		log.Println("Cannot get user:", err)
		writer.Respond(http.StatusUnauthorized, "Credentials are incorrect")
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(received.Password)); err != nil {
		writer.Respond(http.StatusUnauthorized, "Credentials are incorrect")
		return
	}

	writer.Respond(http.StatusOK, struct{
		Message string
		JWT string
	} {
		"Login OK.",
		"123456",
	})
}

func HandlePosts(w http.ResponseWriter, r *http.Request) {
}

func HandlePostByID(w http.ResponseWriter, r *http.Request) {
}
