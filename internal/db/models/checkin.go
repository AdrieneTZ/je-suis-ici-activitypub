package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Checkin struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Content      string    `json:"content"`
	LocationName string    `json:"location_name"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	ActivityID   string    `json:"activity_id"`
	Media        []Media   `json:"media,omitempty"`
	User         *User     `json:"user,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CheckinRepository methods to manipulate checkin data
type CheckinRepository interface {
	CreateCheckin(ctx context.Context, checkin *Checkin) error
	GetCheckinByID(ctx context.Context, id uuid.UUID) (*Checkin, error)
}

// CheckinRepositoryImplement implement functions in checkin repository interface
type CheckinRepositoryImplement struct {
	pool *pgxpool.Pool
}

// NewCheckinRepository create CheckinRepository interface instance
func NewCheckinRepository(pool *pgxpool.Pool) CheckinRepository {
	return &CheckinRepositoryImplement{pool: pool}
}

func (cr *CheckinRepositoryImplement) CreateCheckin(ctx context.Context, checkin *Checkin) error {
	query := `
		INSERT INTO checkins (
			user_id, content, location_name, latitude, longitude, activity_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := cr.pool.QueryRow(ctx, query,
		checkin.UserID, checkin.Content, checkin.LocationName,
		checkin.Latitude, checkin.Longitude, checkin.ActivityID,
	).Scan(&checkin.ID, &checkin.CreatedAt, &checkin.UpdatedAt)

	if err != nil {
		return fmt.Errorf("fail to create checkin: %w", err)
	}

	return nil
}

func (cr *CheckinRepositoryImplement) GetCheckinByID(ctx context.Context, id uuid.UUID) (*Checkin, error) {
	query := `
SELECT
c.id, c.user_id, c.content, c.location_name, c.latitude, c.longitude, 
c.activity_id, c.created_at, c.updated_at,
u.id, u.username, u.display_name, u.avatar_url, u.actor_id
FROM checkins c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1
`

	row := cr.pool.QueryRow(ctx, query, id)

	var checkin Checkin
	var user User

	// get checkin data and user data
	err := row.Scan(
		&checkin.ID, &checkin.UserID, &checkin.Content, &checkin.LocationName, &checkin.Latitude, &checkin.Longitude,
		&checkin.ActivityID, &checkin.CreatedAt, &checkin.UpdatedAt,
		&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.ActorID,
	)

	if err != nil {
		return nil, fmt.Errorf("fail to get checkin by ID: %w", err)
	}

	// get media data
	mediaQuery := `
SELECT id, file_path, file_type, file_size, width, height, created_at
FROM media
WHERE checkin_id = $1
`
	mediaRows, err := cr.pool.Query(ctx, mediaQuery, id)
	if err != nil {
		return nil, fmt.Errorf("fail to query media: %w", err)
	}
	defer mediaRows.Close()

	// if next row exists, mediaRows.Next() return true and scan media data
	// if it's at the end, mediaRows.Next() return false and break the loop
	for mediaRows.Next() {
		var media Media
		err := mediaRows.Scan(
			&media.ID, &media.FilePath, &media.FileType, &media.FileSize,
			&media.Width, &media.Height, &media.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("fail to scan media: %w", err)
		}

		media.CheckinID = id
		checkin.Media = append(checkin.Media, media)
	}

	err = mediaRows.Err()
	if err != nil {
		return nil, fmt.Errorf("error on iterating media rows: %w", err)
	}

	return &checkin, nil
}
