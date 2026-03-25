package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(pool *pgxpool.Pool) error {
	ctx := context.Background()

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		create_at TIMESTAMP DEFAULT NOW(),
		update_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS authors (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS publishers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS books (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(500) NOT NULL,
		blurb TEXT,
		image TEXT,
		goodreads_link TEXT,
		custom_link TEXT,
		create_at TIMESTAMP DEFAULT NOW(),
		update_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS book_authors (
		book_id UUID REFERENCES books(id) ON DELETE CASCADE,
		author_id UUID REFERENCES authors(id) ON DELETE CASCADE,
		PRIMARY KEY (book_id, author_id)
	);

	CREATE TABLE IF NOT EXISTS reading_progress (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		book_id UUID REFERENCES books(id) ON DELETE CASCADE,
		pages INTEGER NOT NULL DEFAULT 0,
		pages_read INTEGER NOT NULL DEFAULT 0,
		status VARCHAR(50) NOT NULL DEFAULT 'reading',
		create_at TIMESTAMP DEFAULT NOW(),
		update_at TIMESTAMP DEFAULT NOW(),
		UNIQUE(user_id, book_id)
	);

	CREATE INDEX IF NOT EXISTS idx_reading_progress_user_id ON reading_progress(user_id);
	CREATE INDEX IF NOT EXISTS idx_reading_progress_book_id ON reading_progress(book_id);
	`

	_, err := pool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
