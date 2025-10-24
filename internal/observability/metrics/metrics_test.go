package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/metrics"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/tests/mock"
)

func TestNewMetricsService_DefaultMetricsEnabled(t *testing.T) {
	cfg := &config.Metrics{EnableDefaultMetrics: true}
	m := metrics.NewMetricsService(cfg)

	// We can’t inspect registry directly (unexported), but we can use Handler() to verify defaults exist.
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	m.Handler().ServeHTTP(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body := rec.Body.String()
	assert.Contains(t, body, "go_gc_duration_seconds")
	assert.Contains(t, body, "process_cpu_seconds_total")
}

func TestNewMetricsService_DefaultMetricsDisabled(t *testing.T) {
	cfg := &config.Metrics{EnableDefaultMetrics: false}
	m := metrics.NewMetricsService(cfg)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	m.Handler().ServeHTTP(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body := rec.Body.String()
	// Should be empty since no default metrics or custom collectors registered.
	assert.Empty(t, strings.TrimSpace(body))
}

func TestRegister_CustomCollector(t *testing.T) {
	cfg := &config.Metrics{}
	mockC := new(mock.MockCollector)
	mockC.On("Register", testifymock.Anything).Return()

	m := metrics.NewMetricsService(cfg)
	m.Register(mockC)

	mockC.AssertCalled(t, "Register", testifymock.Anything)
	assert.True(t, mockC.Registered)
}
