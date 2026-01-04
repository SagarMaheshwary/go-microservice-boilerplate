package service

import (
	"context"
	"sync/atomic"
)

type HealthService interface {
	Check(ctx context.Context) HealthStatus
	SetReady(ready bool)
}

type HealthStatus struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

type DependencyHealthCheck func(ctx context.Context) error

type healthService struct {
	ready  atomic.Bool
	checks map[string]DependencyHealthCheck
}

func NewHealthService(checks map[string]DependencyHealthCheck) HealthService {
	h := &healthService{
		checks: checks,
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

	for name, fn := range h.checks {
		if err := fn(ctx); err != nil {
			status.Status = "unready"
			status.Details[name] = err.Error()
		} else {
			status.Details[name] = "ok"
		}
	}

	return status
}
