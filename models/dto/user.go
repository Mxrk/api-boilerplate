package dto

import (
	"time"
)

// LoginRequest is a struct for the login endpoint.
type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserTokenResponse struct {
	User UserRequest `json:"user,omitempty"`
}

type UserRequest struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
