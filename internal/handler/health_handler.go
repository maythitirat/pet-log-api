package handler

import (
	"net/http"
	"time"

	"github.com/maythitirat/pet-log-api/pkg/response"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Uptime    string  `json:"uptime"`
	Version   string  `json:"version"`
}

// Health returns the health status of the API
// @Summary Health check
// @Description Returns the health status of the API
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(h.startTime).String(),
		Version:   "1.0.0",
	}

	response.JSON(w, http.StatusOK, resp)
}

// Ready returns whether the API is ready to accept requests
// @Summary Readiness check
// @Description Returns whether the API is ready to accept requests
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}
