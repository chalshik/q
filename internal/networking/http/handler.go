package http

import (
	"encoding/json"
	"net/http"
	"time"

	service "dispatcher/internal/services"
)

type createJobRequest struct {
	UserID string `json:"user_id"`
	Prompt string `json:"prompt"`
}

type createJobResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	jobService *service.JobService
}

func NewHandler(js *service.JobService) *Handler {
	return &Handler{
		jobService: js,
	}
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use the service layer to process data persistence and enqueuing
	job, err := h.jobService.Create(r.Context(), req.UserID, req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := createJobResponse{
		ID:        job.ID,
		Status:    job.Status,
		CreatedAt: job.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted represents async work started
	json.NewEncoder(w).Encode(response)
}
