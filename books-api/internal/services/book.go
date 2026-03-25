package service

import (
	"context"
	"errors"
	"fmt"

	"books/internal/cache"
	"books/internal/models"
	"books/internal/repository"

	"github.com/jackc/pgx/v5"
)

type BookRepositoryInterface interface {
	Create(ctx context.Context, book *models.Book, authorNames []string) error
	GetByID(ctx context.Context, id string) (*models.Book, error)
	List(ctx context.Context, limit, offset int) ([]models.Book, error)
	Update(ctx context.Context, book *models.Book) error
}

type BookService struct {
	bookRepo   BookRepositoryInterface
	authorRepo *repository.AuthorRepository
	cache      *cache.Cache
}

func NewBookService(bookRepo *repository.BookRepository, authorRepo *repository.AuthorRepository, cache *cache.Cache) *BookService {
	return &BookService{bookRepo: bookRepo, authorRepo: authorRepo, cache: cache}
}

func NewBookServiceWithRepos(bookRepo BookRepositoryInterface, authorRepo *repository.AuthorRepository, cache *cache.Cache) *BookService {
	return &BookService{bookRepo: bookRepo, authorRepo: authorRepo, cache: cache}
}

type CreateBookRequest struct {
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Blurb         string   `json:"blurb"`
	Image         string   `json:"image"`
	GoodreadsLink string   `json:"goodreads_link"`
	CustomLink    string   `json:"custom_link"`
}

func (s *BookService) Create(ctx context.Context, req CreateBookRequest) (*models.Book, error) {
	book := &models.Book{
		Title:         req.Title,
		Blurb:         req.Blurb,
		Image:         req.Image,
		GoodreadsLink: req.GoodreadsLink,
		CustomLink:    req.CustomLink,
	}

	if err := s.bookRepo.Create(ctx, book, req.Authors); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx)

	return book, nil
}

func (s *BookService) GetByID(ctx context.Context, id string) (*models.Book, error) {
	if s.cache != nil {
		cacheKey := fmt.Sprintf("book:%s", id)
		var cachedBook models.Book
		if err := s.cache.Get(ctx, cacheKey, &cachedBook); err == nil {
			return &cachedBook, nil
		}
	}

	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("book:%s", id)
		_ = s.cache.Set(ctx, cacheKey, book)
	}

	return book, nil
}

func (s *BookService) List(ctx context.Context, limit, offset int) ([]models.Book, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("books:list:%d:%d", limit, offset)
		var cachedBooks []models.Book
		if err := s.cache.Get(ctx, cacheKey, &cachedBooks); err == nil {
			return cachedBooks, nil
		}
	}

	books, err := s.bookRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("books:list:%d:%d", limit, offset)
		_ = s.cache.Set(ctx, cacheKey, books)
	}

	return books, nil
}

type UpdateBookRequest struct {
	Title         *string   `json:"title,omitempty"`
	Authors       *[]string `json:"authors,omitempty"`
	Blurb         *string   `json:"blurb,omitempty"`
	Image         *string   `json:"image,omitempty"`
	GoodreadsLink *string   `json:"goodreads_link,omitempty"`
	CustomLink    *string   `json:"custom_link,omitempty"`
}

func (s *BookService) Update(ctx context.Context, id string, req UpdateBookRequest) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Blurb != nil {
		book.Blurb = *req.Blurb
	}
	if req.Image != nil {
		book.Image = *req.Image
	}
	if req.GoodreadsLink != nil {
		book.GoodreadsLink = *req.GoodreadsLink
	}
	if req.CustomLink != nil {
		book.CustomLink = *req.CustomLink
	}

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx)

	return book, nil
}

func (s *BookService) invalidateCache(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.DeletePattern(ctx, "book:*")
	_ = s.cache.DeletePattern(ctx, "books:*")
}
