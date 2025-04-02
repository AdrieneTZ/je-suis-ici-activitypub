package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Database struct {
	Pool *pgxpool.Pool
}

// Config database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabase create new database connection
func NewDatabase(cfg Config) (*Database, error) {
	// build connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode)

	if cfg.Password != "" {
		connectionString += fmt.Sprintf(" password=%s", cfg.Password)
	}

	// setup connection pool
	connPoolCfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		// error message contain original error with %w
		return nil, fmt.Errorf("unable to parse pool config: %w", err)
	}

	// max connection and min connection, can be set higher if needed
	connPoolCfg.MaxConns = 10
	connPoolCfg.MinConns = 2
	connPoolCfg.MaxConnLifetime = 45 * time.Minute
	connPoolCfg.MaxConnIdleTime = 15 * time.Minute

	// create connection pool
	connPool, err := pgxpool.NewWithConfig(context.Background(), connPoolCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// test connection
	err = connPool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Database{Pool: connPool}, nil
}

// Close close database connection
func (db *Database) Close() {
	// to prevent database connection initial fail
	// if Close is accidentally called more than two times,
	// those calls might cause panic
	if db != nil && db.Pool != nil {
		db.Pool.Close()
		db.Pool = nil
	}
}
