package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistRepository struct {
	pool *pgxpool.Pool
}

func NewPlaylistRepository(pool *pgxpool.Pool) *PlaylistRepository {
	return &PlaylistRepository{pool: pool}
}

func (r *PlaylistRepository) Create(ctx context.Context, playlist *domain.Playlist) error {
	const query = `
		INSERT INTO playlists (
			user_id, title, description, is_private, created_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		playlist.UserID,
		playlist.Title,
		playlist.Description,
		playlist.IsPrivate,
		playlist.CreatedAt,
	).Scan(&playlist.ID)
	if err != nil {
		return fmt.Errorf("create playlist: %w", mapError(err))
	}

	return nil
}

func (r *PlaylistRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Playlist, error) {
	const query = `
		SELECT id, user_id, title, description, is_private, created_at
		FROM playlists
		WHERE id = $1;
	`

	var m playlistModel
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.UserID,
		&m.Title,
		&m.Description,
		&m.IsPrivate,
		&m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get playlist by id: %s: %w", id, mapError(err))
	}

	return m.toDomain(), nil
}

func (r *PlaylistRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Playlist, error) {
	const query = `
		SELECT id, user_id, title, description, is_private, created_at
		FROM playlists
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get playlists by user id %s: %w", userID, mapError(err))
	}
	defer rows.Close()

	var playlists []*domain.Playlist
	for rows.Next() {
		var m playlistModel
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Title,
			&m.Description,
			&m.IsPrivate,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan playlist: %w", err)
		}
		playlists = append(playlists, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlists: %w", err)
	}

	return playlists, nil
}

func (r *PlaylistRepository) Update(ctx context.Context, playlist *domain.Playlist) error {
	const query = `
		UPDATE playlists
		SET title = $1, description = $2, is_private = $3
		WHERE id = $4;
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		playlist.Title,
		playlist.Description,
		playlist.IsPrivate,
		playlist.ID,
	)
	if err != nil {
		return fmt.Errorf("update playlist %s: %w", playlist.ID, mapError(err))
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: playlist id %s", domain.ErrNotFound, playlist.ID)
	}

	return nil
}

func (r *PlaylistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM playlists
		WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete playlist %s: %w", id, mapError(err))
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: playlist id %s", domain.ErrNotFound, id)
	}

	return nil
}

func (r *PlaylistRepository) AddTrack(ctx context.Context, playlistID, trackID uuid.UUID, position int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if position <= 0 {
		const maxPosQuery = `
			SELECT COALESCE(MAX(position), 0) + 1
			FROM playlist_tracks
			WHERE playlist_id = $1;
		`

		if err := tx.QueryRow(ctx, maxPosQuery, playlistID).Scan(&position); err != nil {
			return fmt.Errorf("calculate next track position: %w", mapError(err))
		}
	} else {
		const shiftQuery = `
			UPDATE playlist_tracks
			SET position = position + 1
			WHERE playlist_id = $1 AND position >= $2;
		`

		if _, err := tx.Exec(ctx, shiftQuery, playlistID, position); err != nil {
			return fmt.Errorf("shift playlist tracks: %w", mapError(err))
		}
	}

	const insertQuery = `
		INSERT INTO playlist_tracks (playlist_id, track_id, position, added_at)
		VALUES ($1, $2, $3, $4);
	`

	_, err = tx.Exec(ctx, insertQuery, playlistID, trackID, position, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("insert track to playlist: %w", mapError(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *PlaylistRepository) RemoveTrack(ctx context.Context, playlistID, trackID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const deleteQuery = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = $1 AND track_id = $2
		RETURNING position;
	`

	var deletedPosition int
	err = tx.QueryRow(ctx, deleteQuery, playlistID, trackID).Scan(&deletedPosition)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrTrackNotInPlaylist
		}

		return fmt.Errorf("delete track from playlist: %w", mapError(err))
	}

	const shiftQuery = `
		UPDATE playlist_tracks
		SET position = position - 1
		WHERE playlist_id = $1 AND position > $2;
	`

	if _, err := tx.Exec(ctx, shiftQuery, playlistID, deletedPosition); err != nil {
		return fmt.Errorf("rebalance track positions: %w", mapError(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *PlaylistRepository) ChangeTrackPosition(ctx context.Context, playlistID, trackID uuid.UUID, newPosition int) error {
	if newPosition <= 0 {
		return fmt.Errorf("%w: position must be positive", domain.ErrInvalidInput)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const getPosQuery = `
		SELECT position
		FROM playlist_tracks
		WHERE playlist_id = $1 AND track_id = $2;
	`

	var oldPosition int
	err = tx.QueryRow(ctx, getPosQuery, playlistID, trackID).Scan(&oldPosition)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrTrackNotInPlaylist
		}
		return fmt.Errorf("get current track position: %w", mapError(err))
	}

	if oldPosition == newPosition {
		return nil
	}

	if oldPosition < newPosition {
		const shiftDownQuery = `
			UPDATE playlist_tracks
            SET position = position - 1
            WHERE playlist_id = $1 AND position > $2 AND position <= $3;
		`

		if _, err := tx.Exec(ctx, shiftDownQuery, playlistID, oldPosition, newPosition); err != nil {
			return fmt.Errorf("shift tracks down: %w", mapError(err))
		}
	} else {
		const shiftUpQuery = `
            UPDATE playlist_tracks
            SET position = position + 1
            WHERE playlist_id = $1 AND position >= $2 AND position < $3;
        `
		if _, err := tx.Exec(ctx, shiftUpQuery, playlistID, newPosition, oldPosition); err != nil {
			return fmt.Errorf("shift tracks up: %w", mapError(err))
		}
	}

	const updatePosQuery = `
		UPDATE playlist_tracks
        SET position = $1
        WHERE playlist_id = $2 AND track_id = $3;
	`

	if _, err := tx.Exec(ctx, updatePosQuery, newPosition, playlistID, trackID); err != nil {
		return fmt.Errorf("update track target position: %w", mapError(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *PlaylistRepository) GetTracks(ctx context.Context, playlistID uuid.UUID) ([]*domain.PlaylistTrack, error) {
	const query = `
		SELECT playlist_id, track_id, position, added_at
		FROM playlist_tracks
		WHERE playlist_id = $1
		ORDER BY position ASC;
	`

	rows, err := r.pool.Query(ctx, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("get playlist tracks for %s: %w", playlistID, mapError(err))
	}
	defer rows.Close()

	var tracks []*domain.PlaylistTrack
	for rows.Next() {
		var m playlistTrackModel
		if err := rows.Scan(
			&m.PlaylistID,
			&m.TrackID,
			&m.Position,
			&m.AddedAt,
		); err != nil {
			return nil, fmt.Errorf("scan playlist track: %w", err)
		}

		tracks = append(tracks, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlist tracks: %w", err)
	}

	return tracks, nil
}
