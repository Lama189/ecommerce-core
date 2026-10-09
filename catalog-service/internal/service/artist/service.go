package artist

import (
	"context"
	"errors"
	"fmt"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	repo ArtistRepository
}

func NewService(repo ArtistRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, dto CreateArtistDTO) (*ArtistOutputDTO, error) {
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, domain.ErrArtistAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("check existing artist for user %s: %w", userID, err)
	}

	newArtist, err := domain.NewArtist(userID, dto.Name, dto.Bio)
	if err != nil {
		return nil, fmt.Errorf("create artist: %w", err)
	}

	if err := s.repo.Create(ctx, newArtist); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return nil, domain.ErrArtistAlreadyExists
		}
		return nil, fmt.Errorf("save artist: %w", err)
	}

	return toOutputDTO(newArtist), nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ArtistOutputDTO, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrArtistNotFound
		}
		return nil, fmt.Errorf("get artist by id %s: %w", id, err)
	}

	return toOutputDTO(a), nil
}

func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) (*ArtistOutputDTO, error) {
	a, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrArtistNotFound
		}
		return nil, fmt.Errorf("get artist by user id %s: %w", userID, err)
	}

	return toOutputDTO(a), nil
}

func (s *Service) Update(ctx context.Context, currentUserID, artistID uuid.UUID, dto UpdateArtistDTO) (*ArtistOutputDTO, error) {
	a, err := s.repo.GetByID(ctx, artistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrArtistNotFound
		}
		return nil, fmt.Errorf("find artist for update: %w", err)
	}

	if a.UserID != currentUserID {
		return nil, domain.ErrNotArtistOwner
	}

	if err := a.Update(dto.Name, dto.Bio); err != nil {
		return nil, fmt.Errorf("update artist entity: %w", err)
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, fmt.Errorf("save updated artist: %w", err)
	}

	return toOutputDTO(a), nil
}

func (s *Service) Delete(ctx context.Context, currentUserID, artistID uuid.UUID) error {
	a, err := s.repo.GetByID(ctx, artistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrArtistNotFound
		}
		return fmt.Errorf("find artist for delete: %w", err)
	}

	if a.UserID != currentUserID {
		return domain.ErrNotArtistOwner
	}

	if err := s.repo.Delete(ctx, artistID); err != nil {
		return fmt.Errorf("delete artist: %w", err)
	}

	return nil
}
