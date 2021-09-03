package main

import (
	"errors"
	"encoding/json"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/ztdevelops/go-project/src/helpers/custom_types"
	"github.com/rs/cors"
)

func GetCorsWrapper(allowedHeaders, allowedMethods []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
	})
}

func GetJWTMiddleware() *jwtmiddleware.JWTMiddleware {
	validationKeyGetter := func(token *jwt.Token) (interface{}, error) {
		audience := "api_identifier_goes_here"
		checkedAudience := token.Claims.(jwt.MapClaims).VerifyAudience(audience, false)
		if !checkedAudience {
			return token, errors.New("invalid audience")
		}

		issuer := "domain"
		checkIssuer := token.Claims.(jwt.MapClaims).VerifyIssuer(issuer, false)
		if !checkIssuer {
			return token, errors.New("invalid issuer")
		}
		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	}

	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: validationKeyGetter,
		SigningMethod: jwt.SigningMethodHS256,
	})
} 

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://YOUR_DOMAIN/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = custom_types.Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil 
}