package artist

import (
	"time"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CreateArtistDTO struct {
	Name string
	Bio  string
}

type UpdateArtistDTO struct {
	Name string
	Bio  string
}

type ArtistOutputDTO struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toOutputDTO(a *domain.Artist) *ArtistOutputDTO {
	if a == nil {
		return nil
	}

	return &ArtistOutputDTO{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Bio:       a.Bio,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func toOutputDTOList(artists []*domain.Artist) []*ArtistOutputDTO {
	if artists == nil {
		return nil
	}

	list := make([]*ArtistOutputDTO, 0, len(artists))
	for _, a := range artists {
		list = append(list, toOutputDTO(a))
	}

	return list
}
