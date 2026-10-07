package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/Lama189/ecommerce-core/api-geteway/internal/delivery/http"
)

type Client interface {
	Register(ctx context.Context, req http.RegisterRequest) (*http.UserResponse, error)
	Login(ctx context.Context, req http.LoginRequest) (*http.AuthResponse, error)
	Refresh(ctx context.Context, req http.RefreshRequest) (*http.RefreshResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*http.UserResponse, error)
}
