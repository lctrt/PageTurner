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

type mockAuthService struct {
	registerResp *svc.AuthResponse
	registerErr  error
	loginResp    *svc.AuthResponse
	loginErr     error
}

func (m *mockAuthService) Register(ctx context.Context, req svc.RegisterRequest) (*svc.AuthResponse, error) {
	if m.registerErr != nil {
		return nil, m.registerErr
	}
	return m.registerResp, nil
}

func (m *mockAuthService) Login(ctx context.Context, req svc.LoginRequest) (*svc.AuthResponse, error) {
	if m.loginErr != nil {
		return nil, m.loginErr
	}
	return m.loginResp, nil
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := &mockAuthService{
		registerResp: &svc.AuthResponse{
			Token: "test-token",
			User:  models.User{ID: "user-1", Username: "testuser", Email: "test@example.com"},
		},
	}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": "testuser", "email": "test@example.com", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp svc.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Token != "test-token" {
		t.Errorf("expected token 'test-token', got '%s'", resp.Token)
	}
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	mockService := &mockAuthService{}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	mockService := &mockAuthService{}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": ""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	mockService := &mockAuthService{}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": "testuser", "password": "123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Register_UsernameTaken(t *testing.T) {
	mockService := &mockAuthService{
		registerErr: svc.ErrUsernameTaken,
	}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": "testuser", "email": "test@example.com", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := &mockAuthService{
		loginResp: &svc.AuthResponse{
			Token: "test-token",
			User:  models.User{ID: "user-1", Username: "testuser"},
		},
	}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp svc.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Token != "test-token" {
		t.Errorf("expected token 'test-token', got '%s'", resp.Token)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := &mockAuthService{
		loginErr: errors.New("invalid credentials"),
	}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": "testuser", "password": "wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	mockService := &mockAuthService{}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{"username": ""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	mockService := &mockAuthService{}
	handler := handlers.NewAuthHandlerWithInterface(mockService)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

var _ = time.Time{}
