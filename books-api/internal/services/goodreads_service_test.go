package service_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"books/internal/models"
	service "books/internal/services"
)

type mockBookCreator struct {
	createdBook *models.Book
	createErr   error
}

func (m *mockBookCreator) Create(ctx context.Context, req service.CreateBookRequest) (*models.Book, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.createdBook, nil
}

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestGoodreadsService_ImportFromGoodreads_Success(t *testing.T) {
	htmlBody := `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
{
  "@type": "Book",
  "name": "Test Book",
  "author": {"@type": "Person", "name": "Test Author"},
  "description": "A test description",
  "image": "https://example.com/image.jpg"
}
</script>
</head>
<body></body>
</html>`

	mockCreator := &mockBookCreator{
		createdBook: &models.Book{
			ID:            "book-123",
			Title:         "Test Book",
			Authors:       []models.Author{{ID: "a1", Name: "Test Author"}},
			Blurb:         "A test description",
			Image:         "https://example.com/image.jpg",
			GoodreadsLink: "https://goodreads.com/book/123",
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
		},
	}

	mockHTTP := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(htmlBody)),
		},
	}

	svc := service.NewGoodreadsServiceWithDeps(mockCreator, mockHTTP)

	book, err := svc.ImportFromGoodreads(context.Background(), service.GoodreadsImportRequest{
		URL: "https://goodreads.com/book/123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if book.Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got '%s'", book.Title)
	}
}

func TestGoodreadsService_ImportFromGoodreads_CreateError(t *testing.T) {
	htmlBody := `<!DOCTYPE html><html><head><script type="application/ld+json">{"@type": "Book", "name": "Test Book"}</script></head><body></body></html>`

	mockCreator := &mockBookCreator{
		createErr: errors.New("database error"),
	}

	mockHTTP := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(htmlBody)),
		},
	}

	svc := service.NewGoodreadsServiceWithDeps(mockCreator, mockHTTP)

	_, err := svc.ImportFromGoodreads(context.Background(), service.GoodreadsImportRequest{
		URL: "https://goodreads.com/book/123",
	})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGoodreadsService_ImportFromGoodreads_NoTitle(t *testing.T) {
	htmlBody := `<!DOCTYPE html><html><head></head><body></body></html>`

	mockCreator := &mockBookCreator{}
	mockHTTP := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(htmlBody)),
		},
	}

	svc := service.NewGoodreadsServiceWithDeps(mockCreator, mockHTTP)

	_, err := svc.ImportFromGoodreads(context.Background(), service.GoodreadsImportRequest{
		URL: "https://goodreads.com/book/123",
	})

	if err == nil {
		t.Error("expected error for missing title, got nil")
	}
}

func TestGoodreadsService_ParseGoodreadsPage_HTTPError(t *testing.T) {
	mockCreator := &mockBookCreator{}
	mockHTTP := &mockHTTPClient{
		err: errors.New("network error"),
	}

	svc := service.NewGoodreadsServiceWithDeps(mockCreator, mockHTTP)

	_, err := svc.ParseGoodreadsPage(context.Background(), "https://goodreads.com/book/123")

	if err == nil {
		t.Error("expected error, got nil")
	}
}
