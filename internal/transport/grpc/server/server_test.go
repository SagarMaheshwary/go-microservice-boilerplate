package server_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transport/grpc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// TestNewServer ensures a GRPCServer struct is created with correct dependencies.
func TestNewServer(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)
	cfg := &config.GRPCServer{}

	srv := server.NewServer(&server.Opts{
		Config: cfg,
		Logger: log,
	})

	require.NotNil(t, srv)
	assert.Equal(t, cfg, srv.Config)
	assert.Equal(t, log, srv.Logger)
	assert.NotNil(t, srv.Server)
}

// TestServeListener verifies server can run on a bufconn listener.
func TestServeListener(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	srv := server.NewServer(&server.Opts{
		Config: &config.GRPCServer{},
		Logger: log,
	})

	lis := bufconn.Listen(bufSize)

	go func() {
		_ = srv.ServeListener(lis)
	}()
	defer srv.Server.Stop()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "gRPC server started")
}

// TestServe verifies server can run on a real TCP listener (ephemeral port).
func TestServe(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	srv := server.NewServer(&server.Opts{
		Config: &config.GRPCServer{
			URL: ":0",
		},
		Logger: log,
	})

	// Use :0 to let OS pick a free port
	go func() {
		_ = srv.Serve()
	}()
	defer srv.Server.Stop()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "gRPC server started")
}
