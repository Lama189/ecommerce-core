package http

import (
	"context"

	"github.com/google/uuid"
)

type UserClient interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Refresh(ctx context.Context, req RefreshRequest) (*RefreshResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}
