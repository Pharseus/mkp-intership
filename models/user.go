package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Fullname     string    `json:"fullname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginRequest model untuk request login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse model untuk response login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
