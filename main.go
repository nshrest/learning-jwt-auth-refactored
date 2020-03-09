// Package main creates a server which handles few api endpoints and restricted api endpoint
// which can be be eccessed via jwt token based authentication.
// Learning JWT authentication and its implementation to go webpages.
package main

import (
	"database/sql"
	"learning-jwt-auth-refactored/controllers"
	"learning-jwt-auth-refactored/driver"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

// func init run before the main function
// this loads all env vars also from .env file in root folder.
func init() {
	gotenv.Load()
}

var db *sql.DB

func main() {
	// Connect to db
	db = driver.ConnectDB()

	controller := controllers.Controller{}

	// NewRouter returns a new router instance. (pointer)
	router := mux.NewRouter()

	// A router.HandleFunc registers routes to be matched and dispatches a handler. (pointer to a route struct)
	// route now has various methods and one we use to respond if the req method is get/post called Methods
	// signup route is registered which calls signup handle function when http.Methods is POST
	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleware(controller.ProtectedEndpoint())).Methods("GET")

	// start server on localhost port 8080
	log.Println("Listen on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
