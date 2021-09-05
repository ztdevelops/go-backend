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
		log.Fatalf("error loading .env:", err)
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
		log.Fatalf("error getting auth client:", err)
	}
	idToken := strings.Split(bearerToken, "Bearer ")[1]

	_, err = client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("wrong token")
	}
	return
}

func LogInWithFirebase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	received := custom_types.User{}
	err := json.NewDecoder(r.Body).Decode(&received)
	if err != nil {
		log.Fatalf("failed to decode user:", err)
	}
	log.Println(received)
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", getEnv("API_KEY"))
	if !received.ReturnSecureToken {
		received.ReturnSecureToken = true
	}
	u, err := json.Marshal(received)
	if err != nil {
		log.Fatalf("error marshalling user:", err)
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(u))
	if err != nil {
		log.Fatalf("error querying firebase api:", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("error reading body:", err)
	}

	var result custom_types.UserReponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("error unmarshalling response:", err)
	}

	json.NewEncoder(w).Encode(result)

	// if err := json.NewDecoder(r.Body).Decode(received); err != nil {
	// 	writer.Respond(http.StatusBadRequest, err.Error())
	// 	return
	// }
}

// package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"log"
// 	"net/http"

// 	"github.com/auth0/go-jwt-middleware"
// 	"github.com/form3tech-oss/jwt-go"
// 	"github.com/rs/cors"
// 	"github.com/ztdevelops/go-project/src/helpers/custom_types"
// )

// func GetCorsWrapper(allowedHeaders, allowedMethods []string) *cors.Cors {
// 	return cors.New(cors.Options{
// 		AllowedMethods: allowedMethods,
// 		AllowedHeaders: allowedHeaders,
// 	})
// }

// func RetrieveToken() {

// }

// func GetJWTMiddleware() *jwtmiddleware.JWTMiddleware {
// 	validationKeyGetter := func(token *jwt.Token) (interface{}, error) {
// 		// audience := "localhost:8000/api"
// 		// checkedAudience := token.Claims.(jwt.MapClaims).VerifyAudience(audience, false)
// 		// if !checkedAudience {
// 		// 	return token, errors.New("invalid audience")
// 		// }

// 		issuer := "https://dev-3tmog62q.us.auth0.com/"
// 		checkIssuer := token.Claims.(jwt.MapClaims).VerifyIssuer(issuer, false)
// 		if !checkIssuer {
// 			return token, errors.New("invalid issuer")
// 		}
// 		cert, err := getPemCert(token)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
// 		return result, nil
// 	}

// 	return jwtmiddleware.New(jwtmiddleware.Options{
// 		ValidationKeyGetter: validationKeyGetter,
// 		SigningMethod: jwt.SigningMethodHS256,
// 	})
// }

// func getPemCert(token *jwt.Token) (string, error) {
// 	cert := ""
// 	resp, err := http.Get("https://dev-3tmog62q.us.auth0.com/.well-known/jwks.json")

// 	if err != nil {
// 		return cert, err
// 	}
// 	defer resp.Body.Close()

// 	var jwks = custom_types.Jwks{}
// 	err = json.NewDecoder(resp.Body).Decode(&jwks)

// 	if err != nil {
// 		return cert, err
// 	}

// 	for k, v := range token.Header {
// 		log.Println("key:", k, "val:", v)
// 	}

// 	for k, _ := range jwks.Keys {
// 		if token.Header["kid"] == jwks.Keys[k].Kid {
// 			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
// 		}
// 	}

// 	if cert == "" {
// 		err := errors.New("unable to find appropriate key")
// 		return cert, err
// 	}

// 	return cert, nil
// }
