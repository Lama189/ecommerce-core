package postgres

import (
	"time"

	"github.com/Lama189/soundwave-platform/user-service/internal/domain"
	"github.com/google/uuid"
)

type userModel struct {
	ID           uuid.UUID `db:"id"`
	Phone        string    `db:"phone"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (m *userModel) toDomain() *domain.User {
	return domain.RestoreUser(
		m.ID,
		m.Phone,
		m.PasswordHash,
		domain.Role(m.Role),
		m.CreatedAt,
		m.UpdatedAt,
	)
}
