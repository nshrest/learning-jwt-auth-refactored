package utils

import (
	"encoding/json"
	"learning-jwt-auth-refactored/models"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// func respondWithError creates response with error message as a body
func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

// func responseJSON creates response with user (without password) as a body
// also writes header as "Content-Type": "application/json"
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func GenerateToken(user models.User) (string, error) {
	var err error

	// jwt = header.payload.secret
	// set secret for jwt-token
	secret := os.Getenv("SECRET")

	// generate a new token which takes a signingmethod(algorithm like HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   user.Username,
		"issuer": "pe-info",
	})

	// Test how does the token looks like
	// spew.Dump(token)
	// return "", nil

	// sign the token with secret key which would be a final token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func ComparePasswords(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	if err != nil {
		log.Println("login unsuccessful", err)
		return false
	}
	return true
}
