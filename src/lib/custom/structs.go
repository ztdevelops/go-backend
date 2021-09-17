package custom

import (
	"net/http"
)

type User struct {
	ID                int
	Email             string `json:"email" gorm:"unique;not null"`
	Password          string `json:"password" gorm:"not null"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

// Struct created for querying firebase APIs.
// Firebase only requires Email, Password, and ReturnSecureToken.
type UserForFirebase struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type userResponse struct {
	Status       int    `json:"status"`
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

type UserSignInResponse struct {
	userResponse
	Registered bool `json:"registered"`
}

type UserSignUpResponse struct {
	userResponse
	Kind string `json:"kind"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Response struct {
	Status  int
	Message interface{}
}

/*
	STRUCT EXTENSIONS
*/
type CustomWriter struct {
	http.ResponseWriter
}

type CustomRequest struct {
	*http.Request
}
