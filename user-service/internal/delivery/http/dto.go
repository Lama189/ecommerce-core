package http

import (
	"time"

	"github.com/google/uuid"
)

type registerRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type loginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type authResponse struct {
	User         userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
