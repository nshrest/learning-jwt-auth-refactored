package models

// A JWT struct holds the jwt token for protected endpoint verification
type JWT struct {
	Token string `json:"token"`
}
