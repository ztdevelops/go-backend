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

type UserReponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	Registered   bool   `json:"registered"`
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
