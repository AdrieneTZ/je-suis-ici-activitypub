package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Media struct {
	ID        uuid.UUID `json:"id"`
	CheckinID uuid.UUID `json:"checkin_id,omitempty"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	FileSize  int       `json:"file_size"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	URL       string    `json:"url,omitempty"` // not in database, generate by server
	CreatedAt time.Time `json:"created_at"`
}

// MediaRepository methods to manipulate media data
type MediaRepository interface {
	CreateMedia(ctx context.Context, media *Media) error
	GetMediaByID(ctx context.Context, id uuid.UUID) (*Media, error)
	UpdateMedia(ctx context.Context, media *Media) error
}

// MediaRepositoryImplement
type MediaRepositoryImplement struct {
	pool *pgxpool.Pool
}

// NewMediaRepository
func NewMediaRepository(pool *pgxpool.Pool) MediaRepository {
	return &MediaRepositoryImplement{pool: pool}
}

// CreateMedia store media path and related information to database
func (mr *MediaRepositoryImplement) CreateMedia(ctx context.Context, media *Media) error {
	query := `
		INSERT INTO media (
			checkin_id, file_path, file_type, file_size, width, height
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	// use nil for checkin_id if it's uuid.Nil to properly handle SQL NULL
	var checkinID interface{}
	if media.CheckinID == uuid.Nil {
		checkinID = nil
	} else {
		checkinID = media.CheckinID
	}

	err := mr.pool.QueryRow(ctx, query,
		checkinID, media.FilePath, media.FileType, media.FileSize, media.Width, media.Height,
	).Scan(&media.ID, &media.CreatedAt)

	if err != nil {
		return fmt.Errorf("fail to create media: %w", err)
	}

	return nil
}

// GetMediaByID
func (mr *MediaRepositoryImplement) GetMediaByID(ctx context.Context, id uuid.UUID) (*Media, error) {
	query := `
		SELECT id, checkin_id, file_path, file_type, file_size, width, height, created_at
		FROM media
		WHERE id = $1
	`

	media := &Media{}
	err := mr.pool.QueryRow(ctx, query, id).Scan(
		&media.ID, &media.CheckinID, &media.FilePath, &media.FileType,
		&media.FileSize, &media.Width, &media.Height, &media.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("fail to get media by ID: %w", err)
	}

	return media, nil
}

// UpdateMedia
func (mr *MediaRepositoryImplement) UpdateMedia(ctx context.Context, media *Media) error {
	query := `
		UPDATE media
		SET checkin_id = $1, file_path = $2, file_type = $3, file_size = $4, width = $5, height = $6
		WHERE id = $7
	`

	_, err := mr.pool.Exec(ctx, query,
		media.CheckinID, media.FilePath, media.FileType,
		media.FileSize, media.Width, media.Height, media.ID,
	)

	if err != nil {
		return fmt.Errorf("fail to update media: %w", err)
	}

	return nil
}
