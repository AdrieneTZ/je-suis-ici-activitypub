package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name,omitempty"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	ActorID      string    `json:"actor_id"`
	PrivateKey   string    `json:"-"`
	PublicKey    string    `json:"public_key,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRepository manipulate user data
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
}

// UserRepositoryImplement implement functions in user repository interface
type UserRepositoryImplement struct {
	pool *pgxpool.Pool
}

// NewUserRepository create UserRepository instance
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &UserRepositoryImplement{pool: pool}
}

func (ur *UserRepositoryImplement) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			username, display_name, email, password_hash, avatar_url, actor_id, private_key, public_key
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := ur.pool.QueryRow(ctx, query,
		user.Username, user.DisplayName, user.Email, user.PasswordHash, user.AvatarURL, user.ActorID, user.PrivateKey, user.PublicKey,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("fail to create user: %w", err)
	}

	return nil
}
