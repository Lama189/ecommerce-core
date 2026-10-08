package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ListeningHistory struct {
	ID         int64
	UserID     uuid.UUID
	TrackID    uuid.UUID
	ListenedAt time.Time
	Completed  bool
}

func NewListeningHistory(userID, trackID uuid.UUID, completed bool) (*ListeningHistory, error) {
	now := time.Now().UTC()
	lh := &ListeningHistory{
		UserID:     userID,
		TrackID:    trackID,
		ListenedAt: now,
		Completed:  completed,
	}

	if err := lh.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	return lh, nil
}

func RestoreListeningHistory(id int64, userID, trackID uuid.UUID, listenedAt time.Time, completed bool) *ListeningHistory {
	return &ListeningHistory{
		ID:         id,
		UserID:     userID,
		TrackID:    trackID,
		ListenedAt: listenedAt,
		Completed:  completed,
	}
}

func (lh *ListeningHistory) Validate() error {
	if lh.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}
	if lh.TrackID == uuid.Nil {
		return fmt.Errorf("track ID cannot be empty")
	}

	return nil
}
