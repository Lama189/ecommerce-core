package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewArtist(userID uuid.UUID, name, bio string) (*Artist, error) {
	now := time.Now().UTC()
	a := &Artist{
		UserID:    userID,
		Name:      strings.TrimSpace(name),
		Bio:       strings.TrimSpace(bio),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := a.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	return a, nil
}

func RestoreArtist(id, userID uuid.UUID, name, bio string, createdAt, updatedAt time.Time) *Artist {
	return &Artist{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Bio:       bio,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (a *Artist) Validate() error {
	if a.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	trimmedName := strings.TrimSpace(a.Name)
	if trimmedName == "" {
		return fmt.Errorf("artist name cannot be empty")
	}
	if len(trimmedName) > 255 {
		return fmt.Errorf("artist name exceeds maximum length of 255 characters")
	}

	if len(a.Bio) > 2000 {
		return fmt.Errorf("artist bio exceeds maximum length of 2000 characters")
	}

	return nil
}
