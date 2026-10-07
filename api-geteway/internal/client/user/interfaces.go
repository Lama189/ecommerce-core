package user

import (
	"context"
	"uuid"

	"github.com/Lama189/ecommerce-core/api-geteway/internal/delivery/http"
)

type Client interface {
	Register(ctx context.Context, req http.RegisterRequest) (*http.UserResponse, error)
	Login(ctx context.Context, req http.LoginRequest) (*http.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*http.UserResponse, error)
}
