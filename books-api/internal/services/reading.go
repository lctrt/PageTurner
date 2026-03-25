package service

import (
	"context"
	"errors"

	"books/internal/models"
	"books/internal/repository"

	"github.com/jackc/pgx/v5"
)

type ReadingProgressRepositoryInterface interface {
	Create(ctx context.Context, progress *models.ReadingProgress) error
	GetByUserAndBook(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error)
	Update(ctx context.Context, progress *models.ReadingProgress) error
	GetByUserID(ctx context.Context, userID string) ([]models.ReadingProgress, error)
}

type ReadingService struct {
	progressRepo ReadingProgressRepositoryInterface
	bookRepo     BookRepositoryInterface
}

func NewReadingService(progressRepo *repository.ReadingProgressRepository, bookRepo *repository.BookRepository) *ReadingService {
	return &ReadingService{progressRepo: progressRepo, bookRepo: bookRepo}
}

func NewReadingServiceWithRepos(progressRepo ReadingProgressRepositoryInterface, bookRepo BookRepositoryInterface) *ReadingService {
	return &ReadingService{progressRepo: progressRepo, bookRepo: bookRepo}
}

type UpdateReadingStatusRequest struct {
	Pages     *int                  `json:"pages,omitempty"`
	PagesRead *int                  `json:"pages_read,omitempty"`
	Status    *models.ReadingStatus `json:"status,omitempty"`
}

func (s *ReadingService) GetStatus(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error) {
	progress, err := s.progressRepo.GetByUserAndBook(ctx, userID, bookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("reading progress not found")
		}
		return nil, err
	}
	return progress, nil
}

func (s *ReadingService) StartReading(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error) {
	_, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	progress := &models.ReadingProgress{
		UserID:    userID,
		BookID:    bookID,
		Pages:     0,
		PagesRead: 0,
		Status:    models.StatusReading,
	}

	if err := s.progressRepo.Create(ctx, progress); err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *ReadingService) UpdateStatus(ctx context.Context, userID, bookID string, req UpdateReadingStatusRequest) (*models.ReadingProgress, error) {
	progress, err := s.progressRepo.GetByUserAndBook(ctx, userID, bookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err := s.bookRepo.GetByID(ctx, bookID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, errors.New("book not found")
				}
				return nil, err
			}

			pages := 0
			pagesRead := 0
			status := models.StatusReading
			if req.Pages != nil {
				pages = *req.Pages
			}
			if req.PagesRead != nil {
				pagesRead = *req.PagesRead
			}
			if req.Status != nil {
				status = *req.Status
			}

			progress = &models.ReadingProgress{
				UserID:    userID,
				BookID:    bookID,
				Pages:     pages,
				PagesRead: pagesRead,
				Status:    status,
			}

			if err := s.progressRepo.Create(ctx, progress); err != nil {
				return nil, err
			}
			return progress, nil
		}
		return nil, err
	}

	if req.Pages != nil {
		progress.Pages = *req.Pages
	}
	if req.PagesRead != nil {
		progress.PagesRead = *req.PagesRead
	}
	if req.Status != nil {
		progress.Status = *req.Status
	}

	if err := s.progressRepo.Update(ctx, progress); err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *ReadingService) GetUserBooks(ctx context.Context, userID string) ([]models.ReadingProgress, error) {
	return s.progressRepo.GetByUserID(ctx, userID)
}
