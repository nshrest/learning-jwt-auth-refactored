package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"learning-jwt-auth-refactored/models"
	userRepository "learning-jwt-auth-refactored/repository/user"
	"learning-jwt-auth-refactored/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var users []models.User

type Controller struct{}

// func Signup handler func accepts responsewriter interface and request struct.
// ResponseWriter { Header() Header, Write([]byte) (int, error), WriteHeader(statusCode int) }
// Request { Method string, URL *url.URL, Header Header, Body io.ReadCloser, Form url.Values, ctx context.Context}   ... many more in structs
// but those are the important ones.
func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// w.Write([]byte("Successfully called signup"))
		// extracts user input as username and password from request body, encrypt password and stores it back to database
		// after successfull (unsuccessfull) operation, return back response to client.
		var user models.User
		var error models.Error

		// NewDecoder returns a decoder struct from a request and Decode method converts the json to provided struct
		json.NewDecoder(r.Body).Decode(&user)

		// check if any empty values provided in request from clinet.
		if user.Username == "" {
			log.Println("Error: signup endpoint invoked")
			error.Message = "Username is missing."
			// send Bad Request http.StatusBadRequest 400
			// w.WriteHeader(http.StatusBadRequest)
			// json.NewEncoder(w).Encode(error)
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			log.Println("Error: signup endpoint invoked")
			error.Message = "Password is missing."
			// send Bad Request http.StatusBadRequest 400
			// w.WriteHeader(http.StatusBadRequest)
			// json.NewEncoder(w).Encode(error)
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if error.Message == "" {
			log.Println("signup endpoint invoked")
		}

		// encrypt password to store in db
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)
		// Add user to db

		userRepo := userRepository.UserRepository{}
		user = userRepo.Signup(db, user)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		// set pasword to nil for client response
		user.Password = ""

		// write header and response
		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)

	}
}

// func Login is  handler function which accepts a responsewriter interface and pointer to a response.
func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("successfully called login"))
		var user models.User
		var jwt models.JWT
		var error models.Error

		log.Println("login endpoint invoked")

		// extract user info from request and update user variable
		json.NewDecoder(r.Body).Decode(&user)

		// Check if username and password is not empty
		if user.Username == "" {
			error.Message = "Username is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			log.Println("login unsuccessful")
			return
		}
		if user.Password == "" {
			error.Message = "Password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			log.Println("login unsuccessful")
			return
		}

		// user password from the login request (plain text)
		// to comapare later with hashed password from db
		userpassword := user.Password

		// check if user exists in database

		userRepo := userRepository.UserRepository{}
		user, err := userRepo.Login(db, user)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Username not found"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				log.Println("login unsuccessful")
				return
			} else {
				log.Println("login unsuccessful.. sth wrong terminating...")
				log.Fatal(err)
			}
		}

		// compare hash password to user provided password
		hashedpassword := user.Password
		isValidPassword := utils.ComparePasswords(hashedpassword, []byte(userpassword))
		if isValidPassword {
			// generate token by passing user
			token, err := utils.GenerateToken(user)
			if err != nil {
				log.Fatal(err)
			}

			jwt.Token = token
			// fmt.Println(token)
			// setting client response with JWT token
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			utils.ResponseJSON(w, jwt)
			log.Println("login success")
		} else {
			error.Message = "Invalid Password"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

	}
}

// TokenVerifyMiddleware is a middleware function sits between protected endpoint and protected endpoint handle function.
// This outputs the protected endpoint handle function which calls the protected endpoint after verifiying auth.
// This does by validating token from the header with secret key
// for every protected endpoint just call this tokenverifymiddleware onwards....
func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error

		// find the token part from header Authorization (bearer token)
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			// Parse, validate, and return a token. (When Parse is successful - token's valid field is true, signature field is written)
			// keyFunc will receive the parsed token and should return the key for validating.
			// If everything is kosher, err will be nil
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}

			if token.Valid {
				// ServeHTTP calls next(w, r)  - in our case next is protectedEndpoint(w, r)
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token."
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
