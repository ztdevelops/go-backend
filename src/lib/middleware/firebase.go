package middleware

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/ztdevelops/go-project/src/lib/custom"
	"google.golang.org/api/option"
)

// InitFirebase initialises the connection to Firebase, which
// will be used for auth purposes. Returns a firebase App
// object if done successfully. Else, it will return an error.
func InitFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("credentials.json")
	return firebase.NewApp(context.Background(), nil, opt)
}

// VerifyToken checks the validity of the token.
// If the token is invalid, an error will be returned.
func VerifyToken(a *firebase.App, r *http.Request) (err error) {
	bearerToken := r.Header.Get("Authorization")
	ctx := r.Context()

	client, err := a.Auth(ctx)
	if err != nil {
		log.Println("error getting auth client:", err)
		return
	}
	idToken := strings.Split(bearerToken, "Bearer ")[1]

	_, err = client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return
	}
	return
}

// LoginWithFirebase queries the Firebase servers 
// in an attempt to authenticate the user.
func LoginWithFirebase(user []byte) (*http.Response, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", custom.GetEnv("API_KEY"))
	return http.Post(url, custom.ContentTypeJSON, bytes.NewBuffer(user))
}
