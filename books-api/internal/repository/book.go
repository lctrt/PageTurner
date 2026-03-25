package repository

import (
	"context"

	"books/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool       *pgxpool.Pool
	authorRepo *AuthorRepository
}

func NewBookRepository(pool *pgxpool.Pool, authorRepo *AuthorRepository) *BookRepository {
	return &BookRepository{pool: pool, authorRepo: authorRepo}
}

func (r *BookRepository) Create(ctx context.Context, book *models.Book, authorNames []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO books (title, blurb, image, goodreads_link, custom_link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, create_at, update_at
	`
	err = tx.QueryRow(ctx, query, book.Title, book.Blurb, book.Image, book.GoodreadsLink, book.CustomLink).
		Scan(&book.ID, &book.CreateAt, &book.UpdateAt)
	if err != nil {
		return err
	}

	for _, name := range authorNames {
		author, err := r.authorRepo.GetOrCreate(ctx, name)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)`, book.ID, author.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *BookRepository) GetByID(ctx context.Context, id string) (*models.Book, error) {
	query := `
		SELECT id, title, blurb, image, goodreads_link, custom_link, create_at, update_at
		FROM books WHERE id = $1
	`
	book := &models.Book{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&book.ID, &book.Title, &book.Blurb, &book.Image,
		&book.GoodreadsLink, &book.CustomLink, &book.CreateAt, &book.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	authors, err := r.authorRepo.GetByBookID(ctx, id)
	if err != nil {
		return nil, err
	}
	book.Authors = authors

	return book, nil
}

func (r *BookRepository) List(ctx context.Context, limit, offset int) ([]models.Book, error) {
	query := `
		SELECT id, title, blurb, image, goodreads_link, custom_link, create_at, update_at
		FROM books ORDER BY create_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Blurb, &book.Image,
			&book.GoodreadsLink, &book.CustomLink, &book.CreateAt, &book.UpdateAt)
		if err != nil {
			return nil, err
		}
		authors, err := r.authorRepo.GetByBookID(ctx, book.ID)
		if err != nil {
			return nil, err
		}
		book.Authors = authors
		books = append(books, book)
	}
	return books, rows.Err()
}

func (r *BookRepository) Update(ctx context.Context, book *models.Book) error {
	query := `
		UPDATE books SET title = $2, blurb = $3, image = $4, goodreads_link = $5, custom_link = $6, update_at = NOW()
		WHERE id = $1
		RETURNING update_at
	`
	return r.pool.QueryRow(ctx, query, book.ID, book.Title, book.Blurb, book.Image, book.GoodreadsLink, book.CustomLink).
		Scan(&book.UpdateAt)
}
