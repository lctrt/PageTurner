package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"books/internal/models"
	svc "books/internal/services"
)

type BookServiceInterface interface {
	Create(ctx context.Context, req svc.CreateBookRequest) (*models.Book, error)
	GetByID(ctx context.Context, id string) (*models.Book, error)
	List(ctx context.Context, limit, offset int) ([]models.Book, error)
	Update(ctx context.Context, id string, req svc.UpdateBookRequest) (*models.Book, error)
}

type BookHandler struct {
	bookService BookServiceInterface
}

func NewBookHandler(bookService *svc.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func NewBookHandlerWithInterface(bookService BookServiceInterface) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req svc.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	book, err := h.bookService.GetByID(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "book not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	books, err := h.bookService.List(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req svc.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.Update(r.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "book not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
