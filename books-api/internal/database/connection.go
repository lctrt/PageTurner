package database

import (
	"context"
	"fmt"
	"time"

	"books/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	const maxRetries = 5
	const initialBackoff = 500 * time.Millisecond

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.New(context.Background(), dsn)
		if err == nil {
			if pingErr := pool.Ping(context.Background()); pingErr == nil {
				return pool, nil
			}
			pool.Close()
		}

		if i < maxRetries-1 {
			backoff := initialBackoff * time.Duration(1<<i)
			fmt.Printf("Database connection attempt %d/%d failed: %v. Retrying in %v...\n", i+1, maxRetries, err, backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
