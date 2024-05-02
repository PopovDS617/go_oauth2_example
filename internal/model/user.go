package model

type User struct {
	ID              string
	Name            string
	Picture         string
	Email           string
	IsEmailVerified bool `json:"verified_email"`
}
