package history

import (
	"context"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type ListeningHistoryRepository interface {
	Record(ctx context.Context, history *domain.ListeningHistory) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ListeningHistory, error)
	GetRecentlyPlayedTrackIDs(ctx context.Context, userID uuid.UUID, limit int) ([]uuid.UUID, error)
}

type TrackRepository interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Track, error)
}
