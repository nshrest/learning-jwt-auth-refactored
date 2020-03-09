package controllers

import (
	"learning-jwt-auth-refactored/utils"
	"log"
	"net/http"
)

// func protectedEndpoint is a handler function which accepts responsewriter interface and a pointer to a response type
func (c Controller) ProtectedEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "Successfully accessed protected endpoint")
		log.Println("protected endpoint invoked")
	}
}
