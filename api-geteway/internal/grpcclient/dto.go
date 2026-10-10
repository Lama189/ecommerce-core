package grpcclient

import (
	"time"

	"github.com/google/uuid"
)

type ArtistResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateArtistRequest struct {
	UserID uuid.UUID
	Name   string
	Bio    string
}

type UpdateArtistRequest struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
	Bio    string
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

