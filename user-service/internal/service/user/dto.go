package user

import (
	"time"

	"github.com/Lama189/soundwave-platform/user-service/internal/domain"
	"github.com/google/uuid"
)

type CreateUserDTO struct {
	Phone    string
	Password string
}

type UserOutputDTO struct {
	ID        uuid.UUID
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TokensOutputDTO struct {
	AccessToken  string
	RefreshToken string
}

type UserWithTokensOutputDTO struct {
	User   UserOutputDTO
	Tokens TokensOutputDTO
}

func toOutputDTO(user *domain.User) *UserOutputDTO {
	if user == nil {
		return nil
	}

	return &UserOutputDTO{
		ID:        user.ID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toTokensOutputDTO(accessToken, refreshToken string) *TokensOutputDTO {
	return &TokensOutputDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func toUserWithTokensOutputDTO(
	user *domain.User,
	accessToken string,
	refreshToken string,
) *UserWithTokensOutputDTO {
	if user == nil {
		return nil
	}

	return &UserWithTokensOutputDTO{
		User: *toOutputDTO(user),
		Tokens: TokensOutputDTO{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

type TokenPayload struct {
	UserID uuid.UUID
	Role   string
}
