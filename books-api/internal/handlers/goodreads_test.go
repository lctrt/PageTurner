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

type mockGoodreadsService struct {
	importResp *models.Book
	importErr  error
}

func (m *mockGoodreadsService) ImportFromGoodreads(ctx context.Context, req svc.GoodreadsImportRequest) (*models.Book, error) {
	if m.importErr != nil {
		return nil, m.importErr
	}
	return m.importResp, nil
}

func TestGoodreadsHandler_Import_Success(t *testing.T) {
	mockService := &mockGoodreadsService{
		importResp: &models.Book{
			ID:            "book-123",
			Title:         "Imported Book",
			Authors:       []models.Author{{ID: "a1", Name: "Test Author"}},
			Blurb:         "An imported description",
			Image:         "https://example.com/image.jpg",
			GoodreadsLink: "https://goodreads.com/book/123",
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
		},
	}
	handler := handlers.NewGoodreadsHandlerWithInterface(mockService)

	body := `{"url": "https://goodreads.com/book/123"}`
	req := httptest.NewRequest(http.MethodPost, "/import", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Import(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp models.Book
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Title != "Imported Book" {
		t.Errorf("expected title 'Imported Book', got '%s'", resp.Title)
	}
}

func TestGoodreadsHandler_Import_InvalidBody(t *testing.T) {
	mockService := &mockGoodreadsService{}
	handler := handlers.NewGoodreadsHandlerWithInterface(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/import", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Import(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGoodreadsHandler_Import_MissingURL(t *testing.T) {
	mockService := &mockGoodreadsService{}
	handler := handlers.NewGoodreadsHandlerWithInterface(mockService)

	body := `{"url": ""}`
	req := httptest.NewRequest(http.MethodPost, "/import", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Import(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGoodreadsHandler_Import_ParseError(t *testing.T) {
	mockService := &mockGoodreadsService{
		importErr: svc.ErrFailedToParseGoodreads,
	}
	handler := handlers.NewGoodreadsHandlerWithInterface(mockService)

	body := `{"url": "https://goodreads.com/book/123"}`
	req := httptest.NewRequest(http.MethodPost, "/import", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Import(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGoodreadsHandler_Import_ServiceError(t *testing.T) {
	mockService := &mockGoodreadsService{
		importErr: errors.New("database error"),
	}
	handler := handlers.NewGoodreadsHandlerWithInterface(mockService)

	body := `{"url": "https://goodreads.com/book/123"}`
	req := httptest.NewRequest(http.MethodPost, "/import", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Import(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
