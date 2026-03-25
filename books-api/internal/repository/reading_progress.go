package repository

import (
	"context"

	"books/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadingProgressRepository struct {
	pool *pgxpool.Pool
}

func NewReadingProgressRepository(pool *pgxpool.Pool) *ReadingProgressRepository {
	return &ReadingProgressRepository{pool: pool}
}

func (r *ReadingProgressRepository) Create(ctx context.Context, progress *models.ReadingProgress) error {
	query := `
		INSERT INTO reading_progress (user_id, book_id, pages, pages_read, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, create_at, update_at
	`
	return r.pool.QueryRow(ctx, query,
		progress.UserID, progress.BookID, progress.Pages, progress.PagesRead, progress.Status,
	).Scan(&progress.ID, &progress.CreateAt, &progress.UpdateAt)
}

func (r *ReadingProgressRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error) {
	query := `
		SELECT id, user_id, book_id, pages, pages_read, status, create_at, update_at
		FROM reading_progress WHERE user_id = $1 AND book_id = $2
	`
	progress := &models.ReadingProgress{}
	err := r.pool.QueryRow(ctx, query, userID, bookID).Scan(
		&progress.ID, &progress.UserID, &progress.BookID, &progress.Pages,
		&progress.PagesRead, &progress.Status, &progress.CreateAt, &progress.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *ReadingProgressRepository) Update(ctx context.Context, progress *models.ReadingProgress) error {
	query := `
		UPDATE reading_progress 
		SET pages = $3, pages_read = $4, status = $5, update_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING update_at
	`
	return r.pool.QueryRow(ctx, query,
		progress.ID, progress.UserID, progress.Pages, progress.PagesRead, progress.Status,
	).Scan(&progress.UpdateAt)
}

func (r *ReadingProgressRepository) GetByUserID(ctx context.Context, userID string) ([]models.ReadingProgress, error) {
	query := `
		SELECT id, user_id, book_id, pages, pages_read, status, create_at, update_at
		FROM reading_progress WHERE user_id = $1 ORDER BY update_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []models.ReadingProgress
	for rows.Next() {
		var p models.ReadingProgress
		err := rows.Scan(&p.ID, &p.UserID, &p.BookID, &p.Pages, &p.PagesRead, &p.Status, &p.CreateAt, &p.UpdateAt)
		if err != nil {
			return nil, err
		}
		progress = append(progress, p)
	}
	return progress, rows.Err()
}
