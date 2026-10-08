package postgres

import (
	"context"
	"fmt"

	"github.com/Lama189/soundwave-platform/user-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (phone, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Phone,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("create user: %w", mapError(err))
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `
		SELECT id, phone, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	var m userModel

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Phone,
		&m.PasswordHash,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by id %d: %w", id, mapError(err))
	}

	return m.toDomain(), nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	const query = `
		SELECT id, phone, password_hash, created_at, updated_at
		FROM users
		WHERE phone = $1;
	`

	var m userModel

	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&m.ID,
		&m.Phone,
		&m.PasswordHash,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by phone %s: %w", phone, mapError(err))
	}

	return m.toDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	const query = `
		UPDATE users
		SET phone = $1, password_hash = $2, updated_at = $3
		WHERE id = $4;
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		user.Phone,
		user.PasswordHash,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("update user: %w", mapError(err))
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: id %s", domain.ErrNotFound, user.ID)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM users
		WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: id %s", domain.ErrNotFound, id)
	}

	return nil
}
