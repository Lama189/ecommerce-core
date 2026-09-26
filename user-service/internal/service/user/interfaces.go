package user

import (
	"context"

	"github.com/Lama189/ecommerce-core/user-service/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

type TokenManager interface {
	GenerateTokenPair(payload TokenPayload) (accessToken, refreshToken string, err error)
	VerifyToken(tokenStr string, expectedType string) (*TokenPayload, error)
}

type UserCacheRepository interface {
	Set(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	SetRefreshToken(ctx context.Context, userID uuid.UUID, token string) error
	GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error
}
