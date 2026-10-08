package postgres

import (
	"context"
	"fmt"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListeningHistoryRepository struct {
	pool *pgxpool.Pool
}

func NewListeningHistoryRepository(pool *pgxpool.Pool) *ListeningHistoryRepository {
	return &ListeningHistoryRepository{pool: pool}
}

func (r *ListeningHistoryRepository) Record(ctx context.Context, history *domain.ListeningHistory) error {
	const query = `
		INSERT INTO listening_history (user_id, track_id, listened_at, completed)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		history.UserID,
		history.TrackID,
		history.ListenedAt,
		history.Completed,
	).Scan(&history.ID)
	if err != nil {
		return fmt.Errorf("record listening history: %w", mapError(err))
	}

	return nil
}

func (r *ListeningHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ListeningHistory, error) {
	const query = `
		SELECT id, user_id, track_id, listened_at, completed
		FROM listening_history
		WHERE user_id = $1
		ORDER BY listened_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get listening history for user %s: %w", userID, mapError(err))
	}
	defer rows.Close()

	var historyList []*domain.ListeningHistory
	for rows.Next() {
		var m listeningHistoryModel
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.TrackID,
			&m.ListenedAt,
			&m.Completed,
		); err != nil {
			return nil, fmt.Errorf("scan listening history record: %w", err)
		}
		historyList = append(historyList, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listening history records: %w", err)
	}

	return historyList, nil
}

func (r *ListeningHistoryRepository) GetRecentlyPlayedTrackIDs(ctx context.Context, userID uuid.UUID, limit int) ([]uuid.UUID, error) {
	const query = `
		SELECT track_id
		FROM listening_history
		WHERE user_id = $1
		GROUP BY track_id
		ORDER BY MAX(listened_at) DESC
		LIMIT $2;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recently played track ids for user %s: %w", userID, mapError(err))
	}
	defer rows.Close()

	var trackIDs []uuid.UUID
	for rows.Next() {
		var trackID uuid.UUID
		if err := rows.Scan(&trackID); err != nil {
			return nil, fmt.Errorf("scan recent track id: %w", err)
		}
		trackIDs = append(trackIDs, trackID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent track ids: %w", err)
	}

	return trackIDs, nil
}
