package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"books/internal/models"

	"github.com/astappiev/microdata"
)

type BookCreatorInterface interface {
	Create(ctx context.Context, req CreateBookRequest) (*models.Book, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GoodreadsService struct {
	bookService BookCreatorInterface
	httpClient  HTTPClient
}

func NewGoodreadsService(bookService *BookService) *GoodreadsService {
	return &GoodreadsService{bookService: bookService, httpClient: http.DefaultClient}
}

func NewGoodreadsServiceWithMock(bookService BookCreatorInterface) *GoodreadsService {
	return &GoodreadsService{bookService: bookService, httpClient: http.DefaultClient}
}

func NewGoodreadsServiceWithDeps(bookService BookCreatorInterface, httpClient HTTPClient) *GoodreadsService {
	return &GoodreadsService{bookService: bookService, httpClient: httpClient}
}

type GoodreadsImportRequest struct {
	URL string `json:"url"`
}

type GoodreadsBookData struct {
	Title         string
	Authors       []string
	Blurb         string
	Image         string
	GoodreadsLink string
}

func isBookItem(item *microdata.Item) bool {
	for _, t := range item.Types {
		lower := strings.ToLower(t)
		if strings.Contains(lower, "book") {
			return true
		}
	}
	return false
}

func extractStringValues(values []interface{}) []string {
	var authors []string
	for _, val := range values {
		switch v := val.(type) {
		case string:
			if v != "" {
				authors = append(authors, v)
			}
		case *microdata.Item:
			if name := extractNameFromItem(v); name != "" {
				authors = append(authors, name)
			}
		case []interface{}:
			authors = append(authors, extractStringValues(v)...)
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok && name != "" {
				authors = append(authors, name)
			}
		}
	}
	return authors
}

func extractNameFromItem(item *microdata.Item) string {
	for name, values := range item.Properties {
		lower := strings.ToLower(name)
		if lower == "name" || lower == "givenname" || lower == "familyname" {
			for _, val := range values {
				if v, ok := val.(string); ok && v != "" {
					return v
				}
			}
		}
	}
	for _, itemType := range item.Types {
		lower := strings.ToLower(itemType)
		if lower == "person" || strings.Contains(lower, "author") {
			for name, values := range item.Properties {
				if strings.Contains(strings.ToLower(name), "name") {
					for _, val := range values {
						if v, ok := val.(string); ok && v != "" {
							return v
						}
					}
				}
			}
		}
	}
	return ""
}

func (s *GoodreadsService) ParseGoodreadsPage(ctx context.Context, url string) (*GoodreadsBookData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &GoodreadsBookData{
		GoodreadsLink: url,
	}

	md, err := microdata.ParseHTML(strings.NewReader(string(body)), "text/html", url)
	if err != nil {
		return nil, err
	}

	for _, item := range md.Items {
		if !isBookItem(item) {
			continue
		}

		for name, values := range item.Properties {
			switch strings.ToLower(name) {
			case "name":
				if len(values) > 0 {
					if v, ok := values[0].(string); ok {
						data.Title = v
					}
				}
			case "author":
				data.Authors = append(data.Authors, extractStringValues(values)...)
			case "description":
				if len(values) > 0 {
					if v, ok := values[0].(string); ok {
						data.Blurb = v
					}
				}
			case "image":
				if len(values) > 0 {
					if v, ok := values[0].(string); ok {
						data.Image = v
					}
				}
			}
		}
	}

	return data, nil
}

func (s *GoodreadsService) ImportFromGoodreads(ctx context.Context, req GoodreadsImportRequest) (*models.Book, error) {
	data, err := s.ParseGoodreadsPage(ctx, req.URL)
	if err != nil {
		return nil, err
	}

	if data.Title == "" {
		return nil, ErrFailedToParseGoodreads
	}

	book, err := s.bookService.Create(ctx, CreateBookRequest{
		Title:         data.Title,
		Authors:       data.Authors,
		Blurb:         data.Blurb,
		Image:         data.Image,
		GoodreadsLink: data.GoodreadsLink,
	})
	if err != nil {
		return nil, err
	}

	return book, nil
}

var ErrFailedToParseGoodreads = errors.New("failed to parse goodreads page: could not extract book title")
