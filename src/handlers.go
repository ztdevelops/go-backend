package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
	// "time"

	// jwtmiddleware "github.com/auth0/go-jwt-middleware"
	// "github.com/form3tech-oss/jwt-go"
	"github.com/ztdevelops/go-project/src/helpers/custom_types"
	// "golang.org/x/crypto/bcrypt"
)

// HandleRoutes initialises the connections to all the explicitly coded routes.
func (a *App) HandleRoutes() {
	app, err := InitFirebase()
	if err != nil {
		log.Fatalf("failed to init firebase:", err)
	}

	a.App = *app
	a.Router.HandleFunc("/", DefaultHandler).Methods((custom_types.RGET))
	a.Router.HandleFunc("/test", TestAPIHandler).Methods((custom_types.RGET))

	// APIs (No need for auth)
	a.Router.HandleFunc("/api/signup", NotImplemented).Methods((custom_types.RPOST))
	a.Router.HandleFunc("/api/signin", LogInWithFirebase).Methods((custom_types.RPOST))

	// APIs (Require auth)
	a.Router.HandleFunc("/api/test", a.TestVerifyToken).Methods(custom_types.RGET)
	// a.Router.Handle("/api/test", handleWithJWT(jwt, NotImplemented)).Methods(custom_types.RGET)
	http.Handle("/", a.Router)
}

func (a *App) TestVerifyToken(w http.ResponseWriter, r *http.Request) {
	if err := VerifyToken(&a.App, r); err != nil {
		log.Println("failed to verify token:", err)
		return
	}
	log.Println("token ok.")
}

// func handleWithJWT(middleware *jwtmiddleware.JWTMiddleware, f func(http.ResponseWriter, *http.Request)) http.Handler {
// 	httpHandler := http.HandlerFunc(f)
// 	return middleware.Handler(httpHandler)
// }

// func generateToken(userID int) (token string, err error) {
// 	os.Setenv("ACCESS_SECRET", "secret")
// 	atClaims := jwt.MapClaims{}
// 	atClaims["authorized"] = true
// 	atClaims["user_id"] = userID
// 	atClaims["exp"] = time.Now().Add(time.Minute * 10).Unix()
// 	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
// 	token, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
// 	return
// }

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
		// {Username: "user 1", Password: "123"},
		// {Username: "user 2", Password: "456"},
	}
	json.NewEncoder(w).Encode(users)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented.")
}

// HandleSignUp creates an account based on the submitted request.
// func HandleSignUp(w http.ResponseWriter, r *http.Request) {
// 	log.Println(custom_types.ENDPOINT_HIT, "sign up")
// 	writer := Writer{w}
// 	writer.SetContentType(custom_types.ContentTypeJSON)

// 	received := &custom_types.User{}
// 	err := json.NewDecoder(r.Body).Decode(received)
// 	if err != nil {
// 		log.Println("Cannot decode:", err)
// 		writer.Respond(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if received.Username == "" || received.Password == "" {
// 		log.Println("Received empty credentials. Rejecting request.")
// 		writer.Respond(http.StatusBadRequest, "Failed to create account: invalid credentials")
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(received.Password), 8)
// 	if err != nil {
// 		log.Println("Error hashing password:", err)
// 		writer.Respond(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	submitted := &custom_types.User{
// 		Username: received.Username,
// 		Password: string(hashedPassword),
// 	}

// 	if err := SharedApp.SignUp(*submitted); err != nil {
// 		log.Println("Cannot sign up:", err)
// 		writer.Respond(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	writer.Respond(http.StatusOK, "User successfully created.")
// }

// HandleSignIn attempts to log the user into their accounts based on the submitted credentials.
// func HandleSignIn(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.RemoteAddr, custom_types.ENDPOINT_HIT, "sign in")
// 	writer := Writer{w}
// 	writer.SetContentType(custom_types.ContentTypeJSON)
	
// 	received := &custom_types.User{}
// 	if err := json.NewDecoder(r.Body).Decode(received); err != nil {
// 		writer.Respond(http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	user, err := SharedApp.GetUser(received.Username)
// 	if err != nil {
// 		log.Println("Cannot get user:", err)
// 		writer.Respond(http.StatusUnauthorized, "Credentials are incorrect")
// 		return
// 	}
// 	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(received.Password)); err != nil {
// 		writer.Respond(http.StatusUnauthorized, "Credentials are incorrect")
// 		return
// 	}

// 	token, err := generateToken(user.ID)
// 	if err != nil {
// 		writer.Respond(http.StatusInternalServerError, "Unable to generate token")
// 		return
// 	}

// 	writer.Respond(http.StatusOK, struct{
// 		Message string
// 		JWT string
// 	} {
// 		"Login OK.",
// 		token,
// 	})
// }
