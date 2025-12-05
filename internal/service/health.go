package service

import (
	"context"
	"sync/atomic"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
)

type HealthService interface {
	Check(ctx context.Context) HealthStatus
	SetReady(ready bool)
}

type HealthStatus struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

type healthService struct {
	database database.DatabaseService
	cache    cache.CacheService
	ready    atomic.Bool
}

type HealthServiceOpts struct {
	Database database.DatabaseService
	Cache    cache.CacheService
}

func NewHealthService(opts *HealthServiceOpts) HealthService {
	h := &healthService{
		database: opts.Database,
		cache:    opts.Cache,
	}
	h.ready.Store(true)
	return h
}

func (h *healthService) SetReady(ready bool) {
	h.ready.Store(ready)
}

func (h *healthService) Check(ctx context.Context) HealthStatus {
	if !h.ready.Load() {
		return HealthStatus{
			Status:  "unready",
			Details: map[string]string{"service": "shutting down"},
		}
	}

	status := HealthStatus{Status: "ready", Details: map[string]string{}}

	if err := h.database.Ping(ctx); err != nil {
		status.Status = "unready"
		status.Details["database"] = err.Error()
	} else {
		status.Details["database"] = "ok"
	}

	if err := h.cache.Ping(ctx); err != nil {
		status.Status = "unready"
		status.Details["cache"] = err.Error()
	} else {
		status.Details["cache"] = "ok"
	}

	return status
}
