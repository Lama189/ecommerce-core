package history

import (
	"context"
	"fmt"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/Lama189/soundwave-platform/catalog-service/internal/service/track"
	"github.com/google/uuid"
)

type Service struct {
	history ListeningHistoryRepository
	tracks  TrackRepository
}

func NewService(history ListeningHistoryRepository, tracks TrackRepository) *Service {
	return &Service{
		history: history,
		tracks:  tracks,
	}
}

func (s *Service) RecordPlayback(ctx context.Context, userID, trackID uuid.UUID, completed bool) error {
	record, err := domain.NewListeningHistory(userID, trackID, completed)
	if err != nil {
		return fmt.Errorf("create listening history: %w", err)
	}

	if err := s.history.Record(ctx, record); err != nil {
		return fmt.Errorf("record history: %w", err)
	}

	return nil
}

func (s *Service) GetHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*HistoryItemDTO, error) {
	historyRecords, err := s.history.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get history for user %s: %w", userID, err)
	}

	if len(historyRecords) == 0 {
		return []*HistoryItemDTO{}, nil
	}

	trackIDs := make([]uuid.UUID, 0, len(historyRecords))
	for _, h := range historyRecords {
		trackIDs = append(trackIDs, h.TrackID)
	}

	tracks, err := s.tracks.GetByIDs(ctx, trackIDs)
	if err != nil {
		return nil, fmt.Errorf("get tracks for history: %w", err)
	}

	trackMap := make(map[uuid.UUID]*domain.Track, len(tracks))
	for _, t := range tracks {
		trackMap[t.ID] = t
	}

	items := make([]*HistoryItemDTO, 0, len(historyRecords))
	for _, h := range historyRecords {
		t, ok := trackMap[h.TrackID]
		if !ok {
			continue
		}
		items = append(items, &HistoryItemDTO{
			ID:         h.ID,
			UserID:     h.UserID,
			TrackID:    h.TrackID,
			ListenedAt: h.ListenedAt,
			Completed:  h.Completed,
			TrackTitle: t.Title,
			ArtistID:   t.ArtistID,
			Duration:   t.Duration,
			CoverKey:   t.CoverKey,
		})
	}

	return items, nil
}

func (s *Service) GetRecentlyPlayed(ctx context.Context, userID uuid.UUID, limit int) ([]*track.TrackOutputDTO, error) {
	trackIDs, err := s.history.GetRecentlyPlayedTrackIDs(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recently played track ids: %w", err)
	}

	if len(trackIDs) == 0 {
		return []*track.TrackOutputDTO{}, nil
	}

	tracks, err := s.tracks.GetByIDs(ctx, trackIDs)
	if err != nil {
		return nil, fmt.Errorf("get tracks by ids: %w", err)
	}

	trackMap := make(map[uuid.UUID]*domain.Track, len(tracks))
	for _, t := range tracks {
		trackMap[t.ID] = t
	}

	result := make([]*track.TrackOutputDTO, 0, len(trackIDs))
	for _, id := range trackIDs {
		if t, ok := trackMap[id]; ok {
			result = append(result, &track.TrackOutputDTO{
				ID:         t.ID,
				ArtistID:   t.ArtistID,
				AlbumID:    t.AlbumID,
				Title:      t.Title,
				Duration:   t.Duration,
				CoverKey:   t.CoverKey,
				FileKey:    t.FileKey,
				PreviewKey: t.PreviewKey,
				Status:     string(t.Status),
				CreatedAt:  t.CreatedAt,
				UpdatedAt:  t.UpdatedAt,
			})
		}
	}

	return result, nil
}
