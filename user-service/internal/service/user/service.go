package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Lama189/soundwave-platform/user-service/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	repo   UserRepository
	cache  UserCacheRepository
	hasher PasswordHasher
	jwt    TokenManager
}

func NewService(repo UserRepository, cache UserCacheRepository, hasher PasswordHasher, jwt TokenManager) *Service {
	return &Service{
		repo:   repo,
		cache:  cache,
		hasher: hasher,
		jwt:    jwt,
	}
}

func (s *Service) Create(ctx context.Context, dto CreateUserDTO) (*UserOutputDTO, error) {
	passwordHash, err := s.hasher.Hash(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUser, err := domain.NewUser(dto.Phone, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("create User: %w", err)
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	return toOutputDTO(newUser), nil
}

func (s *Service) Login(ctx context.Context, phone, password string) (*UserWithTokensOutputDTO, error) {
	u, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	if !s.hasher.Compare(password, u.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	payload := TokenPayload{
		UserID: u.ID,
		Role:   string(u.Role),
	}

	accessToken, refreshToken, err := s.jwt.GenerateTokenPair(payload)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := s.cache.SetRefreshToken(ctx, u.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("save session to cache: %w", err)
	}

	return toUserWithTokensOutputDTO(u, accessToken, refreshToken), nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokensOutputDTO, error) {
	payload, err := s.jwt.VerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, fmt.Errorf("verify refresh token: %w", err)
	}

	savedToken, err := s.cache.GetRefreshToken(ctx, payload.UserID)
	if err != nil || savedToken != refreshToken {
		return nil, ErrTokenRevoked
	}

	u, err := s.repo.GetByID(ctx, payload.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user for token refresh: %w", err)
	}

	newPayload := TokenPayload{
		UserID: u.ID,
		Role:   string(u.Role),
	}

	newAccessToken, newRefreshToken, err := s.jwt.GenerateTokenPair(newPayload)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := s.cache.SetRefreshToken(ctx, u.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("save session to cache: %w", err)
	}

	return toTokensOutputDTO(newAccessToken, newRefreshToken), nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*UserOutputDTO, error) {
	cachedUser, err := s.cache.Get(ctx, id)
	if err == nil && cachedUser != nil {
		return toOutputDTO(cachedUser), nil
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	_ = s.cache.Set(ctx, u)

	return toOutputDTO(u), nil
}
