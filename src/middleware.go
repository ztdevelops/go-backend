package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/ztdevelops/go-project/src/helpers/custom_types"
	"google.golang.org/api/option"
)

func GetCorsWrapper(allowedHeaders, allowedMethods []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
	})
}

func getEnv(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("error loading .env:", err)
	}
	log.Println(key, os.Getenv(key))
	return os.Getenv(key)
}

func InitFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("credentials.json")
	return firebase.NewApp(context.Background(), nil, opt)
}

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

func LogInWithFirebase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	received := custom_types.User{}
	err := json.NewDecoder(r.Body).Decode(&received)
	if err != nil {
		log.Println("failed to decode user:", err)
	}

	// Default to always true, as we want the secure token for authentication purposes.
	received.ReturnSecureToken = true
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", getEnv("API_KEY"))
	u, err := json.Marshal(received)
	if err != nil {
		log.Println("error marshalling user:", err)
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(u))
	if err != nil {
		log.Println("error querying firebase api:", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading body:", err)
	}

	var result custom_types.UserReponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("error unmarshalling response:", err)
	}

	json.NewEncoder(w).Encode(result)

	// if err := json.NewDecoder(r.Body).Decode(received); err != nil {
	// 	writer.Respond(http.StatusBadRequest, err.Error())
	// 	return
	// }
}