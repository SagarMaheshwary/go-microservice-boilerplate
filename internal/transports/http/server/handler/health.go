package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
)

type Opts struct {
	HealthService service.HealthService
	Logger        logger.Logger
}

type HealthHandler struct {
	healthService service.HealthService
	logger        logger.Logger
}

func NewHealthHandler(opts *Opts) *HealthHandler {
	return &HealthHandler{healthService: opts.HealthService, logger: opts.Logger}
}

func (h *HealthHandler) Livez(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	writeJSON(w, map[string]string{"status": "ok"}, h.logger)
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	status := h.healthService.Check(r.Context())
	w.Header().Set("Content-Type", "application/json")

	if status.Status == "ready" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	writeJSON(w, status, h.logger)
}

func writeJSON(w http.ResponseWriter, data any, log logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("failed to write JSON: " + err.Error())
	}
}
