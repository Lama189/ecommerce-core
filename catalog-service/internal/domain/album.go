package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID
	ArtistID    uuid.UUID
	Title       string
	Description string
	CoverKey    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAlbum(artistID uuid.UUID, title, description, coverKey string) (*Album, error) {
	now := time.Now().UTC()
	a := &Album{
		ArtistID:    artistID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		CoverKey:    strings.TrimSpace(coverKey),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := a.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	return a, nil
}

func RestoreAlbum(
	id, artistID uuid.UUID,
	title, description, coverKey string,
	createdAt, updatedAt time.Time,
) *Album {
	return &Album{
		ID:          id,
		ArtistID:    artistID,
		Title:       title,
		Description: description,
		CoverKey:    coverKey,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (a *Album) Validate() error {
	if a.ArtistID == uuid.Nil {
		return fmt.Errorf("artist ID cannot be empty")
	}

	trimmedTitle := strings.TrimSpace(a.Title)
	if trimmedTitle == "" {
		return fmt.Errorf("album title cannot be empty")
	}
	if len(trimmedTitle) > 255 {
		return fmt.Errorf("album title exceeds maximum length of 255 characters")
	}

	if len(a.Description) > 2000 {
		return fmt.Errorf("album description exceeds maximum length of 2000 characters")
	}

	if len(a.CoverKey) > 512 {
		return fmt.Errorf("cover key exceeds maximum length of 512 characters")
	}

	return nil
}

func (a *Album) Update(title, description, coverKey string) error {
	a.Title = strings.TrimSpace(title)
	a.Description = strings.TrimSpace(description)
	if trimmedCoverKey := strings.TrimSpace(coverKey); trimmedCoverKey != "" {
		a.CoverKey = trimmedCoverKey
	}
	a.UpdatedAt = time.Now().UTC()

	return a.Validate()
}
