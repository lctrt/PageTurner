package repository

import (
	"context"

	"books/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepository struct {
	pool *pgxpool.Pool
}

func NewAuthorRepository(pool *pgxpool.Pool) *AuthorRepository {
	return &AuthorRepository{pool: pool}
}

func (r *AuthorRepository) Create(ctx context.Context, author *models.Author) error {
	query := `INSERT INTO authors (name) VALUES ($1) RETURNING id`
	return r.pool.QueryRow(ctx, query, author.Name).Scan(&author.ID)
}

func (r *AuthorRepository) GetByID(ctx context.Context, id string) (*models.Author, error) {
	query := `SELECT id, name FROM authors WHERE id = $1`
	author := &models.Author{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&author.ID, &author.Name)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (r *AuthorRepository) GetOrCreate(ctx context.Context, name string) (*models.Author, error) {
	author := &models.Author{}
	err := r.pool.QueryRow(ctx, `SELECT id, name FROM authors WHERE name = $1`, name).Scan(&author.ID, &author.Name)
	if err == nil {
		return author, nil
	}

	query := `INSERT INTO authors (name) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id, name`
	err = r.pool.QueryRow(ctx, query, name).Scan(&author.ID, &author.Name)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (r *AuthorRepository) GetByBookID(ctx context.Context, bookID string) ([]models.Author, error) {
	query := `
		SELECT a.id, a.name FROM authors a
		JOIN book_authors ba ON ba.author_id = a.id
		WHERE ba.book_id = $1
	`
	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var author models.Author
		if err := rows.Scan(&author.ID, &author.Name); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, rows.Err()
}
