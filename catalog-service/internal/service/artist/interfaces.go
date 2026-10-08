package artist

import (
	"context"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type ArtistRepository interface {
	Create(ctx context.Context, artist *domain.Artist) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Artist, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Artist, error)
	Update(ctx context.Context, artist *domain.Artist) error
	Delete(ctx context.Context, id uuid.UUID) error
}
