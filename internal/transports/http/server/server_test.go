package server_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/metrics"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/tests/mock"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestNewServer ensures a HttpServer struct is created with correct dependencies.
func TestNewServer(t *testing.T) {
	log := logger.NewZerologLogger("info", io.Discard)
	cfg := &config.HTTPServer{URL: ":0"}

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	srv := server.NewServer(&server.Opts{
		Config:   cfg,
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
	})

	require.NotNil(t, srv)
	assert.Equal(t, cfg, srv.Config)
	assert.Equal(t, log, srv.Logger)
	assert.NotNil(t, srv.Server)
	assert.NotNil(t, srv.Server.Handler)
}

// TestServeListener verifies server can run on a custom listener.
func TestServeListener(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	srv := server.NewServer(&server.Opts{
		Config:   &config.HTTPServer{},
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
	})

	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer lis.Close()

	go func() {
		_ = srv.ServeListener(lis)
	}()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "HTTP server started")

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = srv.Server.Shutdown(ctx)
}

// TestServeListener_Metrics verifies metrics endpoint works when metricsService
// is passed in NewServer()
func TestServeListener_Metrics(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	metricsService := metrics.NewMetricsService(
		&config.Metrics{EnableDefaultMetrics: false},
	)

	srv := server.NewServer(&server.Opts{
		Config:   &config.HTTPServer{},
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
		Metrics:  metricsService,
	})

	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer lis.Close()

	go func() {
		_ = srv.ServeListener(lis)
	}()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "HTTP server started")

	// Check that /metrics returns 200
	resp, err := http.Get("http://" + lis.Addr().String() + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = srv.Server.Shutdown(ctx)
}

// TestServe verifies server can run on a real TCP listener (ephemeral port).
func TestServe(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)
	cfg := &config.HTTPServer{URL: ":0"}

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	srv := server.NewServer(&server.Opts{
		Config:   cfg,
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
	})

	go func() {
		_ = srv.Serve()
	}()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "HTTP server started")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = srv.Server.Shutdown(ctx)
}
