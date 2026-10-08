package history

import (
	"time"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type RecordPlaybackDTO struct {
	TrackID   uuid.UUID
	Completed bool
}

type HistoryOutputDTO struct {
	ID         int64
	UserID     uuid.UUID
	TrackID    uuid.UUID
	ListenedAt time.Time
	Completed  bool
}

type HistoryItemDTO struct {
	ID         int64
	UserID     uuid.UUID
	TrackID    uuid.UUID
	ListenedAt time.Time
	Completed  bool
	TrackTitle string
	ArtistID   uuid.UUID
	Duration   time.Duration
	CoverKey   string
}

func toOutputDTO(h *domain.ListeningHistory) *HistoryOutputDTO {
	if h == nil {
		return nil
	}

	return &HistoryOutputDTO{
		ID:         h.ID,
		UserID:     h.UserID,
		TrackID:    h.TrackID,
		ListenedAt: h.ListenedAt,
		Completed:  h.Completed,
	}
}

func toOutputDTOList(historyList []*domain.ListeningHistory) []*HistoryOutputDTO {
	if historyList == nil {
		return nil
	}

	list := make([]*HistoryOutputDTO, 0, len(historyList))
	for _, h := range historyList {
		list = append(list, toOutputDTO(h))
	}

	return list
}
