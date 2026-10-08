package v1

import (
	"context"

	"github.com/Lama189/soundwave-platform/user-service/internal/service/user"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, dto user.CreateUserDTO) (*user.UserOutputDTO, error)
	Login(ctx context.Context, phone, password string) (*user.UserWithTokensOutputDTO, error)
	Refresh(ctx context.Context, refreshToken string) (*user.TokensOutputDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.UserOutputDTO, error)
}
