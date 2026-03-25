package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"books/internal/models"
	svc "books/internal/services"
)

type GoodreadsServiceInterface interface {
	ImportFromGoodreads(ctx context.Context, req svc.GoodreadsImportRequest) (*models.Book, error)
}

type GoodreadsHandler struct {
	goodreadsService GoodreadsServiceInterface
}

func NewGoodreadsHandler(goodreadsService *svc.GoodreadsService) *GoodreadsHandler {
	return &GoodreadsHandler{goodreadsService: goodreadsService}
}

func NewGoodreadsHandlerWithInterface(goodreadsService GoodreadsServiceInterface) *GoodreadsHandler {
	return &GoodreadsHandler{goodreadsService: goodreadsService}
}

func (h *GoodreadsHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req svc.GoodreadsImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	book, err := h.goodreadsService.ImportFromGoodreads(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == svc.ErrFailedToParseGoodreads {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
