package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TrackLike struct {
	UserID    uuid.UUID
	TrackID   uuid.UUID
	CreatedAt time.Time
}

func NewTrackLike(userID, trackID uuid.UUID) (*TrackLike, error) {
	now := time.Now().UTC()
	l := &TrackLike{
		UserID:    userID,
		TrackID:   trackID,
		CreatedAt: now,
	}

	if err := l.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	return l, nil
}

func RestoreTrackLike(userID, trackID uuid.UUID, createdAt time.Time) *TrackLike {
	return &TrackLike{
		UserID:    userID,
		TrackID:   trackID,
		CreatedAt: createdAt,
	}
}

func (l *TrackLike) Validate() error {
	if l.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}
	if l.TrackID == uuid.Nil {
		return fmt.Errorf("track ID cannot be empty")
	}

	return nil
}
