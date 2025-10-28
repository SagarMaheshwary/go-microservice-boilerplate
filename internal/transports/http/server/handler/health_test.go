package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/tests/mock"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/http/server/handler"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

// TestLivezHandler_OK verifies that /livez always returns 200 and {"status": "ok"}.
func TestLivezHandler_OK(t *testing.T) {
	log := logger.NewZerologLogger("info", io.Discard)
	mockSvc := new(mock.MockHealthService)

	h := handler.NewHealthHandler(&handler.Opts{
		HealthService: mockSvc,
		Logger:        log,
	})

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	rec := httptest.NewRecorder()

	h.Livez(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

// TestReadyzHandler_Ready verifies that /readyz returns 200 when status is "ready".
func TestReadyzHandler_Ready(t *testing.T) {
	log := logger.NewZerologLogger("info", io.Discard)

	mockSvc := new(mock.MockHealthService)
	mockSvc.On("Check", testifymock.Anything).
		Return(service.HealthStatus{Status: "ready"}).
		Once()

	h := handler.NewHealthHandler(&handler.Opts{
		HealthService: mockSvc,
		Logger:        log,
	})

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	h.Readyz(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp service.HealthStatus
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ready", resp.Status)

	mockSvc.AssertExpectations(t)
}

// TestReadyzHandler_NotReady verifies that /readyz returns 503 when status != "ready".
func TestReadyzHandler_NotReady(t *testing.T) {
	log := logger.NewZerologLogger("info", io.Discard)

	mockSvc := new(mock.MockHealthService)
	mockSvc.On("Check", testifymock.Anything).
		Return(service.HealthStatus{Status: "unready", Details: map[string]string{"db": "down"}}).
		Once()

	h := handler.NewHealthHandler(&handler.Opts{
		HealthService: mockSvc,
		Logger:        log,
	})

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	h.Readyz(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp service.HealthStatus
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "unready", resp.Status)
	assert.Equal(t, "down", resp.Details["db"])

	mockSvc.AssertExpectations(t)
}
