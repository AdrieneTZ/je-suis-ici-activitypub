package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ExecuteMigrations execute migration when start running the app
func ExecuteMigrations(cfg Config) error {
	// build database connection string
	connectionString := fmt.Sprintf(
		"cockroachdb://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// create migration instance
	m, err := migrate.New("file://internal/db/migrations", connectionString)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// execute migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Database migrations execute completely!")
	return nil
}
