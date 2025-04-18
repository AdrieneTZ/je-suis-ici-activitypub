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
	GetCheckinByActivityID(ctx context.Context, activityID string) (*Checkin, error)
	GetCheckinsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Checkin, error)
	GetGlobalFeed(ctx context.Context, limit, offest int) ([]Checkin, error)
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

func (cr *CheckinRepositoryImplement) GetCheckinByActivityID(ctx context.Context, activityID string) (*Checkin, error) {
	query := `
		SELECT c.id, c.user_id, c.content, c.location_name, c.latitude, c.longitude,
		c.activity_id, c.created_at, c.updated_at,
FROM checkins c
WHERE activity_id = $1
`

	row := cr.pool.QueryRow(ctx, query, activityID)

	var checkin Checkin

	err := row.Scan(
		&checkin.ID, &checkin.UserID, &checkin.Content, &checkin.LocationName, &checkin.Latitude, &checkin.Longitude,
		&checkin.ActivityID, &checkin.CreatedAt, &checkin.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("fail to get checkin by activity ID: %w", err)
	}

	return &checkin, nil
}

func (cr *CheckinRepositoryImplement) GetCheckinsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Checkin, error) {
	query := `
		SELECT c.id, c.user_id, c.content, c.location_name, c.latitude, c.longitude,
		c.activity_id, c.created_at, c.updated_at,
		u.id, u.username, u.display_name, u.avatar_url, u.actor_id
		FROM checkins c
		JOIN users u ON c.user_id = u.id
		WHERE c.user_id = $1
	ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := cr.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("fail to get checkins by user ID: %w", err)
	}
	defer rows.Close()

	var checkins []Checkin

	for rows.Next() {
		var checkin Checkin
		var user User

		err := rows.Scan(
			&checkin.ID, &checkin.UserID, &checkin.Content, &checkin.LocationName, &checkin.Latitude, &checkin.Longitude,
			&checkin.ActivityID, &checkin.CreatedAt, &checkin.UpdatedAt,
			&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.ActorID,
		)

		if err != nil {
			return nil, fmt.Errorf("fail to scan checkin: %w", err)
		}

		checkin.User = &user
		checkins = append(checkins, checkin)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error on iterating checkin rows: %w", err)
	}

	// get each checkin's media data
	for i := range checkins {
		mediaQuery := `
			SELECT id, file_path, file_type, file_size, width, height, created_at
			FROM media
			WHERE checkin_id = $1
		`

		mediaRows, err := cr.pool.Query(ctx, mediaQuery, checkins[i].ID)
		if err != nil {
			return nil, fmt.Errorf("fail to query media: %w", err)
		}

		for mediaRows.Next() {
			var media Media
			err := mediaRows.Scan(
				&media.ID, &media.FilePath, &media.FileType, &media.FileSize,
				&media.Width, &media.Height, &media.CreatedAt,
			)

			if err != nil {
				mediaRows.Close()
				return nil, fmt.Errorf("fail to scan media: %w", err)
			}

			media.CheckinID = checkins[i].ID
			checkins[i].Media = append(checkins[i].Media, media)
		}

		if err := mediaRows.Err(); err != nil {
			mediaRows.Close()
			return nil, fmt.Errorf("error on iterating media rows: %w", err)
		}

		mediaRows.Close()
	}

	return checkins, nil
}

func (cr *CheckinRepositoryImplement) GetGlobalFeed(ctx context.Context, limit, offest int) ([]Checkin, error) {
	query := `
		SELECT c.id, c.user_id, c.content, c.location_name, c.latitude, c.longitude,
			c.activity_id, c.created_at, c.updated_at,
			u.id, u.username, u.display_name, u.avatar_url, u.actor_id
		FROM checkins c
		JOIN users u ON c.user_id = u.id
		ORDER BY c.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := cr.pool.Query(ctx, query, limit, offest)
	if err != nil {
		return nil, fmt.Errorf("fail to get global feed: %w", err)
	}
	defer rows.Close()

	var checkins []Checkin

	for rows.Next() {
		var checkin Checkin
		var user User

		err := rows.Scan(
			&checkin.ID, &checkin.UserID, &checkin.Content, &checkin.LocationName, &checkin.Latitude, &checkin.Longitude,
			&checkin.ActivityID, &checkin.CreatedAt, &checkin.UpdatedAt,
			&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.ActorID,
		)

		if err != nil {
			return nil, fmt.Errorf("fail to scan checkin: %w", err)
		}

		checkin.User = &user
		checkins = append(checkins, checkin)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating checkin rows: %w", err)
	}

	// get each checkin's media data
	for i := range checkins {
		mediaQuery := `
			SELECT id, file_path, file_type, file_size, width, height, created_at
			FROM media
			WHERE checkin_id = $1
		`

		mediaRows, err := cr.pool.Query(ctx, mediaQuery, checkins[i].ID)
		if err != nil {
			return nil, fmt.Errorf("fail to get media: %w", err)
		}

		for mediaRows.Next() {
			var media Media

			err := mediaRows.Scan(
				&media.ID, &media.FilePath, &media.FileType, &media.FileSize,
				&media.Width, &media.Height, &media.CreatedAt,
			)
			if err != nil {
				mediaRows.Close()
				return nil, fmt.Errorf("fail to scan media: %w", err)
			}

			media.CheckinID = checkins[i].ID
			checkins[i].Media = append(checkins[i].Media, media)
		}

		mediaRows.Close()

		err = mediaRows.Err()
		if err != nil {
			return nil, fmt.Errorf("error iterating media rows: %w", err)
		}
	}

	return checkins, nil
}
