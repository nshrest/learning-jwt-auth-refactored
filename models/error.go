package models

// A Error struct holds the custom error message which will be sent to clients back if needed.
type Error struct {
	Message string `json:"message"`
}
