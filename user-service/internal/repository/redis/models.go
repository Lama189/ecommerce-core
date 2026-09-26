package redis

import (
	"time"

	"github.com/Lama189/ecommerce-core/user-service/internal/domain"
	"github.com/google/uuid"
)

type userCacheModel struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func fromDomain(u *domain.User) userCacheModel {
	return userCacheModel{
		ID:        u.ID,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (m userCacheModel) toDomain() *domain.User {
	return domain.RestoreUser(m.ID, m.Phone, "", m.CreatedAt, m.UpdatedAt)
}
