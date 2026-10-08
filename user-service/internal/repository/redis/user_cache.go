package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Lama189/soundwave-platform/user-service/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserCacheRepository struct {
	client          *redis.Client
	userTTL         time.Duration
	refreshTokenTTL time.Duration
}

func NewCacheRepository(client *redis.Client, userTTL time.Duration, refreshTokenTTL time.Duration) *UserCacheRepository {
	return &UserCacheRepository{
		client:          client,
		userTTL:         userTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (r *UserCacheRepository) Set(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%s", user.ID)

	data, err := json.Marshal(fromDomain(user))
	if err != nil {
		return fmt.Errorf("marshal user cache: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.userTTL).Err(); err != nil {
		return fmt.Errorf("redis set user: %w", err)
	}

	return nil
}

func (r *UserCacheRepository) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", id)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		return nil, fmt.Errorf("redis get user: %w", err)
	}

	var m userCacheModel
	if err := json.Unmarshal(data, &m); err != nil {
		_ = r.Delete(ctx, id)
		return nil, fmt.Errorf("unmarshal user cache: %w", err)
	}

	return m.toDomain(), nil
}

func (r *UserCacheRepository) Delete(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("user:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *UserCacheRepository) SetRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	key := fmt.Sprintf("refresh_token:%s", userID)

	if err := r.client.Set(ctx, key, token, r.refreshTokenTTL).Err(); err != nil {
		return fmt.Errorf("redis set refresh token: %w", err)
	}

	return nil
}

func (r *UserCacheRepository) GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("refresh_token:%s", userID)

	token, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCacheMiss
		}

		return "", fmt.Errorf("redis get refresh token: %w", err)
	}

	return token, nil
}

func (r *UserCacheRepository) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("refresh_token:%s", userID)
	return r.client.Del(ctx, key).Err()
}
