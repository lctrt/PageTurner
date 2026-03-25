package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"books/internal/models"
	service "books/internal/services"
)

type mockBookRepository struct {
	books      map[string]*models.Book
	createErr  error
	getByIDErr error
	listErr    error
	updateErr  error
	listResult []models.Book
}

func newMockBookRepository() *mockBookRepository {
	return &mockBookRepository{
		books: make(map[string]*models.Book),
	}
}

func (m *mockBookRepository) Create(ctx context.Context, book *models.Book, authorNames []string) error {
	if m.createErr != nil {
		return m.createErr
	}
	book.ID = "book-123"
	book.CreateAt = time.Now()
	book.UpdateAt = time.Now()
	book.Authors = make([]models.Author, len(authorNames))
	for i, name := range authorNames {
		book.Authors[i] = models.Author{ID: "author-" + name, Name: name}
	}
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepository) GetByID(ctx context.Context, id string) (*models.Book, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	book, ok := m.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (m *mockBookRepository) List(ctx context.Context, limit, offset int) ([]models.Book, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	if m.listResult != nil {
		return m.listResult, nil
	}
	result := make([]models.Book, 0)
	for _, book := range m.books {
		result = append(result, *book)
	}
	return result, nil
}

func (m *mockBookRepository) Update(ctx context.Context, book *models.Book) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	book.UpdateAt = time.Now()
	m.books[book.ID] = book
	return nil
}

func TestBookService_Create_Success(t *testing.T) {
	mockRepo := newMockBookRepository()
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	book, err := bookSvc.Create(context.Background(), service.CreateBookRequest{
		Title:   "Test Book",
		Authors: []string{"Author One", "Author Two"},
		Blurb:   "A test book",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if book.Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got '%s'", book.Title)
	}

	if len(book.Authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(book.Authors))
	}
}

func TestBookService_GetByID_Success(t *testing.T) {
	mockRepo := newMockBookRepository()
	mockRepo.books["book-1"] = &models.Book{ID: "book-1", Title: "Existing Book"}
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	book, err := bookSvc.GetByID(context.Background(), "book-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if book.Title != "Existing Book" {
		t.Errorf("expected title 'Existing Book', got '%s'", book.Title)
	}
}

func TestBookService_GetByID_NotFound(t *testing.T) {
	mockRepo := newMockBookRepository()
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	_, err := bookSvc.GetByID(context.Background(), "nonexistent")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBookService_List_Success(t *testing.T) {
	mockRepo := newMockBookRepository()
	mockRepo.listResult = []models.Book{
		{ID: "book-1", Title: "Book One"},
		{ID: "book-2", Title: "Book Two"},
	}
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	books, err := bookSvc.List(context.Background(), 10, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}
}

type capturingBookRepository struct {
	mockBookRepository
	capturedLimit  int
	capturedOffset int
}

func (c *capturingBookRepository) List(ctx context.Context, limit, offset int) ([]models.Book, error) {
	c.capturedLimit = limit
	c.capturedOffset = offset
	return c.mockBookRepository.List(ctx, limit, offset)
}

func TestBookService_List_DefaultLimit(t *testing.T) {
	mockRepo := &capturingBookRepository{
		mockBookRepository: mockBookRepository{listResult: []models.Book{}},
	}
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	_, _ = bookSvc.List(context.Background(), 0, -1)

	if mockRepo.capturedLimit != 20 {
		t.Errorf("expected default limit 20, got %d", mockRepo.capturedLimit)
	}
	if mockRepo.capturedOffset != 0 {
		t.Errorf("expected default offset 0, got %d", mockRepo.capturedOffset)
	}
}

func TestBookService_Update_Success(t *testing.T) {
	mockRepo := newMockBookRepository()
	mockRepo.books["book-1"] = &models.Book{ID: "book-1", Title: "Original Title"}
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	newTitle := "Updated Title"
	book, err := bookSvc.Update(context.Background(), "book-1", service.UpdateBookRequest{
		Title: &newTitle,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if book.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", book.Title)
	}
}

func TestBookService_Update_NotFound(t *testing.T) {
	mockRepo := newMockBookRepository()
	bookSvc := service.NewBookServiceWithRepos(mockRepo, nil, nil)

	newTitle := "Updated Title"
	_, err := bookSvc.Update(context.Background(), "nonexistent", service.UpdateBookRequest{
		Title: &newTitle,
	})

	if err == nil {
		t.Error("expected error, got nil")
	}
}
