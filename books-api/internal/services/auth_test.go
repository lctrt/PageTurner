package service_test

import (
	"context"
	"errors"
	"testing"

	"books/internal/config"
	"books/internal/models"
	service "books/internal/services"
)

type mockUserRepository struct {
	users      map[string]*models.User
	usernames  map[string]*models.User
	createErr  error
	getByIDErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:     make(map[string]*models.User),
		usernames: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	m.usernames[user.Username] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.users[id], nil
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.usernames[username], nil
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	resp, err := authService.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("expected token, got empty string")
	}

	if resp.User.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", resp.User.Username)
	}

	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", resp.User.Email)
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	_, err := authService.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error on first register, got %v", err)
	}

	_, err = authService.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test2@example.com",
		Password: "password456",
	})

	if err != service.ErrUsernameTaken {
		t.Errorf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	_, err := authService.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	resp, err := authService.Login(context.Background(), service.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("expected token, got empty string")
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	_, err := authService.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	_, err = authService.Login(context.Background(), service.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})

	if err != service.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	_, err := authService.Login(context.Background(), service.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	})

	if err != service.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	resp, err := authService.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	userID, err := authService.ValidateToken(resp.Token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if userID != resp.User.ID {
		t.Errorf("expected userID '%s', got '%s'", resp.User.ID, userID)
	}
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := config.JWTConfig{Secret: "test-secret", Expiration: 24}

	authService := service.NewAuthServiceWithRepo(mockRepo, cfg)

	_, err := authService.ValidateToken("invalid-token")
	if !errors.Is(err, service.ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestAuthService_ValidateToken_WrongSecret(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg1 := config.JWTConfig{Secret: "secret-1", Expiration: 24}
	cfg2 := config.JWTConfig{Secret: "secret-2", Expiration: 24}

	authService1 := service.NewAuthServiceWithRepo(mockRepo, cfg1)

	resp, err := authService1.Register(context.Background(), service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	authService2 := service.NewAuthServiceWithRepo(mockRepo, cfg2)

	_, err = authService2.ValidateToken(resp.Token)
	if !errors.Is(err, service.ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
