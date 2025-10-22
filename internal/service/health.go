package service

import (
	"context"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"gorm.io/gorm"
)

type HealthService interface {
	Check(ctx context.Context) HealthStatus
}

type HealthStatus struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

type healthService struct {
	database *gorm.DB
	cache    cache.CacheService
}

type HealthServiceOpts struct {
	Database database.DatabaseService
	Cache    cache.CacheService
}

func NewHealthService(opts *HealthServiceOpts) HealthService {
	return &healthService{
		database: opts.Database.DB(),
		cache:    opts.Cache,
	}
}

func (h *healthService) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{Status: "ready", Details: map[string]string{}}

	if err := h.database.Exec("SELECT 1").Error; err != nil {
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
