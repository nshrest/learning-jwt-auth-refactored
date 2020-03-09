package models

// A User struct holds the information about users used for signup/login
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
