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
	"books/internal/middleware"
	"books/internal/models"
	svc "books/internal/services"
)

type mockReadingService struct {
	getStatusResp    *models.ReadingProgress
	getStatusErr     error
	updateStatusResp *models.ReadingProgress
	updateStatusErr  error
	getUserBooksResp []models.ReadingProgress
	getUserBooksErr  error
}

func (m *mockReadingService) GetStatus(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error) {
	if m.getStatusErr != nil {
		return nil, m.getStatusErr
	}
	return m.getStatusResp, nil
}

func (m *mockReadingService) UpdateStatus(ctx context.Context, userID, bookID string, req svc.UpdateReadingStatusRequest) (*models.ReadingProgress, error) {
	if m.updateStatusErr != nil {
		return nil, m.updateStatusErr
	}
	return m.updateStatusResp, nil
}

func (m *mockReadingService) GetUserBooks(ctx context.Context, userID string) ([]models.ReadingProgress, error) {
	if m.getUserBooksErr != nil {
		return nil, m.getUserBooksErr
	}
	return m.getUserBooksResp, nil
}

func createTestRequestWithUserID(userID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/reading/book/1", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func TestReadingHandler_GetStatus_Success(t *testing.T) {
	mockService := &mockReadingService{
		getStatusResp: &models.ReadingProgress{
			ID:        "progress-1",
			UserID:    "user-1",
			BookID:    "book-1",
			Pages:     300,
			PagesRead: 100,
			Status:    models.StatusReading,
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		},
	}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	req := createTestRequestWithUserID("user-1")
	req = httptest.NewRequest(http.MethodGet, "/reading/book/book-1", nil)
	w := httptest.NewRecorder()

	handler.GetStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReadingHandler_GetStatus_NotFound(t *testing.T) {
	mockService := &mockReadingService{
		getStatusErr: errors.New("reading progress not found"),
	}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/reading/book/book-1", nil)
	w := httptest.NewRecorder()

	handler.GetStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestReadingHandler_UpdateStatus_Success(t *testing.T) {
	mockService := &mockReadingService{
		updateStatusResp: &models.ReadingProgress{
			ID:        "progress-1",
			UserID:    "user-1",
			BookID:    "book-1",
			Pages:     300,
			PagesRead: 150,
			Status:    models.StatusReading,
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		},
	}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	body := `{"pages_read": 150}`
	req := httptest.NewRequest(http.MethodPut, "/reading/book/book-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp models.ReadingProgress
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.PagesRead != 150 {
		t.Errorf("expected pages read 150, got %d", resp.PagesRead)
	}
}

func TestReadingHandler_UpdateStatus_InvalidBody(t *testing.T) {
	mockService := &mockReadingService{}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPut, "/reading/book/book-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReadingHandler_UpdateStatus_InvalidStatus(t *testing.T) {
	mockService := &mockReadingService{}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	body := `{"status": "invalid_status"}`
	req := httptest.NewRequest(http.MethodPut, "/reading/book/book-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReadingHandler_UpdateStatus_ValidStatuses(t *testing.T) {
	testStatuses := []string{"reading", "finished", "paused"}

	for _, status := range testStatuses {
		t.Run(status, func(t *testing.T) {
			mockService := &mockReadingService{
				updateStatusResp: &models.ReadingProgress{
					ID:     "progress-1",
					Status: models.ReadingStatus(status),
				},
			}
			handler := handlers.NewReadingHandlerWithInterface(mockService)

			body := `{"status": "` + status + `"}`
			req := httptest.NewRequest(http.MethodPut, "/reading/book/book-1", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d for '%s', got %d", http.StatusOK, status, w.Code)
			}
		})
	}
}

func TestReadingHandler_UpdateStatus_ServiceError(t *testing.T) {
	mockService := &mockReadingService{
		updateStatusErr: errors.New("database error"),
	}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	body := `{"pages_read": 150}`
	req := httptest.NewRequest(http.MethodPut, "/reading/book/book-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateStatus(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestReadingHandler_GetUserBooks_Success(t *testing.T) {
	mockService := &mockReadingService{
		getUserBooksResp: []models.ReadingProgress{
			{ID: "progress-1", UserID: "user-1", BookID: "book-1", Status: models.StatusReading},
			{ID: "progress-2", UserID: "user-1", BookID: "book-2", Status: models.StatusFinished},
		},
	}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/reading/my-books", nil)
	w := httptest.NewRecorder()

	handler.GetUserBooks(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp []models.ReadingProgress
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("expected 2 progress records, got %d", len(resp))
	}
}

func TestReadingHandler_GetUserBooks_ServiceError(t *testing.T) {
	mockService := &mockReadingService{
		getUserBooksErr: errors.New("database error"),
	}
	handler := handlers.NewReadingHandlerWithInterface(mockService)

	req := httptest.NewRequest(http.MethodGet, "/reading/my-books", nil)
	w := httptest.NewRecorder()

	handler.GetUserBooks(w, req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1")))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

var _ = createTestRequestWithUserID
