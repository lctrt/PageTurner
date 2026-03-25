package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"books/internal/handlers"
	"books/internal/models"
	svc "books/internal/services"
)

type mockBookService struct {
	createResp *models.Book
	createErr  error
	getResp    *models.Book
	getErr     error
	listResp   []models.Book
	listErr    error
	updateResp *models.Book
	updateErr  error
}

func (m *mockBookService) Create(ctx context.Context, req svc.CreateBookRequest) (*models.Book, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.createResp, nil
}

func (m *mockBookService) GetByID(ctx context.Context, id string) (*models.Book, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.getResp, nil
}

func (m *mockBookService) List(ctx context.Context, limit, offset int) ([]models.Book, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResp, nil
}

func (m *mockBookService) Update(ctx context.Context, id string, req svc.UpdateBookRequest) (*models.Book, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.updateResp, nil
}

func TestBookHandler_Create_Success(t *testing.T) {
	mockService := &mockBookService{
		createResp: &models.Book{
			ID:            "book-1",
			Title:         "Test Book",
			Authors:       []models.Author{{ID: "a1", Name: "Test Author"}},
			Blurb:         "A test description",
			Image:         "https://example.com/image.jpg",
			GoodreadsLink: "https://goodreads.com/book/123",
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
		},
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	body := `{"title": "Test Book", "authors": ["Test Author"], "blurb": "A test description"}`
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp models.Book
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got '%s'", resp.Title)
	}
}

func TestBookHandler_Create_InvalidBody(t *testing.T) {
	mockService := &mockBookService{}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBookHandler_Create_ServiceError(t *testing.T) {
	mockService := &mockBookService{
		createErr: errors.New("database error"),
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	body := `{"title": "Test Book"}`
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestBookHandler_Get_Success(t *testing.T) {
	mockService := &mockBookService{
		getResp: &models.Book{
			ID:    "book-1",
			Title: "Test Book",
		},
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/books/book-1", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req.WithContext(context.WithValue(req.Context(), "test", "test")))

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBookHandler_Get_NotFound(t *testing.T) {
	mockService := &mockBookService{
		getErr: errors.New("book not found"),
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/books/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestBookHandler_List_Success(t *testing.T) {
	mockService := &mockBookService{
		listResp: []models.Book{
			{ID: "book-1", Title: "Book One"},
			{ID: "book-2", Title: "Book Two"},
		},
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/books?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp []models.Book
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("expected 2 books, got %d", len(resp))
	}
}

func TestBookHandler_List_ServiceError(t *testing.T) {
	mockService := &mockBookService{
		listErr: errors.New("database error"),
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestBookHandler_Update_Success(t *testing.T) {
	mockService := &mockBookService{
		updateResp: &models.Book{
			ID:    "book-1",
			Title: "Updated Title",
		},
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	body := `{"title": "Updated Title"}`
	req := httptest.NewRequest(http.MethodPatch, "/books/book-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBookHandler_Update_NotFound(t *testing.T) {
	mockService := &mockBookService{
		updateErr: errors.New("book not found"),
	}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	body := `{"title": "Updated Title"}`
	req := httptest.NewRequest(http.MethodPatch, "/books/nonexistent", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestBookHandler_Update_InvalidBody(t *testing.T) {
	mockService := &mockBookService{}
	handler := handlers.NewBookHandlerWithInterface(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPatch, "/books/book-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
