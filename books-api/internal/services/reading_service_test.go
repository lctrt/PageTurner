package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"books/internal/models"
	service "books/internal/services"

	"github.com/jackc/pgx/v5"
)

type mockProgressRepository struct {
	progress    map[string]*models.ReadingProgress
	createErr   error
	getErr      error
	updateErr   error
	getByUserID []models.ReadingProgress
}

func newMockProgressRepository() *mockProgressRepository {
	return &mockProgressRepository{
		progress: make(map[string]*models.ReadingProgress),
	}
}

func (m *mockProgressRepository) key(userID, bookID string) string {
	return userID + ":" + bookID
}

func (m *mockProgressRepository) Create(ctx context.Context, progress *models.ReadingProgress) error {
	if m.createErr != nil {
		return m.createErr
	}
	progress.ID = "progress-123"
	progress.CreateAt = time.Now()
	progress.UpdateAt = time.Now()
	m.progress[m.key(progress.UserID, progress.BookID)] = progress
	return nil
}

func (m *mockProgressRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	p, ok := m.progress[m.key(userID, bookID)]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	return p, nil
}

func (m *mockProgressRepository) Update(ctx context.Context, progress *models.ReadingProgress) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	progress.UpdateAt = time.Now()
	m.progress[m.key(progress.UserID, progress.BookID)] = progress
	return nil
}

func (m *mockProgressRepository) GetByUserID(ctx context.Context, userID string) ([]models.ReadingProgress, error) {
	if m.getByUserID != nil {
		return m.getByUserID, nil
	}
	var result []models.ReadingProgress
	for _, p := range m.progress {
		if p.UserID == userID {
			result = append(result, *p)
		}
	}
	return result, nil
}

type mockProgressBookRepository struct {
	books map[string]*models.Book
}

func newMockProgressBookRepository() *mockProgressBookRepository {
	return &mockProgressBookRepository{
		books: make(map[string]*models.Book),
	}
}

func (m *mockProgressBookRepository) Create(ctx context.Context, book *models.Book, authorNames []string) error {
	return nil
}

func (m *mockProgressBookRepository) GetByID(ctx context.Context, id string) (*models.Book, error) {
	book, ok := m.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (m *mockProgressBookRepository) List(ctx context.Context, limit, offset int) ([]models.Book, error) {
	var result []models.Book
	for _, book := range m.books {
		result = append(result, *book)
	}
	return result, nil
}

func (m *mockProgressBookRepository) Update(ctx context.Context, book *models.Book) error {
	m.books[book.ID] = book
	return nil
}

func TestReadingService_GetStatus_Success(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockProgress.progress["user1:book1"] = &models.ReadingProgress{
		ID:        "progress-1",
		UserID:    "user1",
		BookID:    "book1",
		Pages:     300,
		PagesRead: 100,
		Status:    models.StatusReading,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}
	mockBook := newMockProgressBookRepository()
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	status, err := svc.GetStatus(context.Background(), "user1", "book1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status.PagesRead != 100 {
		t.Errorf("expected pages read 100, got %d", status.PagesRead)
	}
}

func TestReadingService_GetStatus_NotFound(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockBook := newMockProgressBookRepository()
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	_, err := svc.GetStatus(context.Background(), "user1", "book1")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestReadingService_StartReading_Success(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockBook := newMockProgressBookRepository()
	mockBook.books["book1"] = &models.Book{ID: "book1", Title: "Test Book"}
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	progress, err := svc.StartReading(context.Background(), "user1", "book1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if progress.Status != models.StatusReading {
		t.Errorf("expected status reading, got %s", progress.Status)
	}

	if progress.Pages != 0 {
		t.Errorf("expected pages 0, got %d", progress.Pages)
	}
}

func TestReadingService_StartReading_BookNotFound(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockBook := newMockProgressBookRepository()
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	_, err := svc.StartReading(context.Background(), "user1", "nonexistent")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestReadingService_UpdateStatus_CreateNew(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockBook := newMockProgressBookRepository()
	mockBook.books["book1"] = &models.Book{ID: "book1", Title: "Test Book"}
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	pages := 300
	pagesRead := 50
	status := models.StatusReading
	progress, err := svc.UpdateStatus(context.Background(), "user1", "book1", service.UpdateReadingStatusRequest{
		Pages:     &pages,
		PagesRead: &pagesRead,
		Status:    &status,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if progress.Pages != 300 {
		t.Errorf("expected pages 300, got %d", progress.Pages)
	}

	if progress.PagesRead != 50 {
		t.Errorf("expected pages read 50, got %d", progress.PagesRead)
	}
}

func TestReadingService_UpdateStatus_UpdateExisting(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockProgress.progress["user1:book1"] = &models.ReadingProgress{
		ID:        "progress-1",
		UserID:    "user1",
		BookID:    "book1",
		Pages:     300,
		PagesRead: 50,
		Status:    models.StatusReading,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}
	mockBook := newMockProgressBookRepository()
	mockBook.books["book1"] = &models.Book{ID: "book1", Title: "Test Book"}
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	pagesRead := 150
	progress, err := svc.UpdateStatus(context.Background(), "user1", "book1", service.UpdateReadingStatusRequest{
		PagesRead: &pagesRead,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if progress.PagesRead != 150 {
		t.Errorf("expected pages read 150, got %d", progress.PagesRead)
	}
}

func TestReadingService_UpdateStatus_BookNotFound(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockBook := newMockProgressBookRepository()
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	_, err := svc.UpdateStatus(context.Background(), "user1", "nonexistent", service.UpdateReadingStatusRequest{})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestReadingService_GetUserBooks_Success(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockProgress.progress["user1:book1"] = &models.ReadingProgress{
		ID:     "progress-1",
		UserID: "user1",
		BookID: "book1",
		Status: models.StatusReading,
	}
	mockProgress.progress["user1:book2"] = &models.ReadingProgress{
		ID:     "progress-2",
		UserID: "user1",
		BookID: "book2",
		Status: models.StatusFinished,
	}
	mockBook := newMockProgressBookRepository()
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	books, err := svc.GetUserBooks(context.Background(), "user1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}
}

func TestReadingService_GetUserBooks_Empty(t *testing.T) {
	mockProgress := newMockProgressRepository()
	mockBook := newMockProgressBookRepository()
	svc := service.NewReadingServiceWithRepos(mockProgress, mockBook)

	books, err := svc.GetUserBooks(context.Background(), "user1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(books) != 0 {
		t.Errorf("expected 0 books, got %d", len(books))
	}
}
