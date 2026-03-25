package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"books/internal/middleware"
	"books/internal/models"
	svc "books/internal/services"
)

type ReadingServiceInterface interface {
	GetStatus(ctx context.Context, userID, bookID string) (*models.ReadingProgress, error)
	UpdateStatus(ctx context.Context, userID, bookID string, req svc.UpdateReadingStatusRequest) (*models.ReadingProgress, error)
	GetUserBooks(ctx context.Context, userID string) ([]models.ReadingProgress, error)
}

type ReadingHandler struct {
	readingService ReadingServiceInterface
}

func NewReadingHandler(readingService *svc.ReadingService) *ReadingHandler {
	return &ReadingHandler{readingService: readingService}
}

func NewReadingHandlerWithInterface(readingService ReadingServiceInterface) *ReadingHandler {
	return &ReadingHandler{readingService: readingService}
}

func (h *ReadingHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	bookID := r.PathValue("bookId")

	progress, err := h.readingService.GetStatus(r.Context(), userID, bookID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "reading progress not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

func (h *ReadingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	bookID := r.PathValue("bookId")

	var req svc.UpdateReadingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != nil {
		validStatuses := map[models.ReadingStatus]bool{
			models.StatusReading:  true,
			models.StatusFinished: true,
			models.StatusPaused:   true,
		}
		if !validStatuses[*req.Status] {
			http.Error(w, "invalid status value", http.StatusBadRequest)
			return
		}
	}

	progress, err := h.readingService.UpdateStatus(r.Context(), userID, bookID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

func (h *ReadingHandler) GetUserBooks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	progress, err := h.readingService.GetUserBooks(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
