package config_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearEnv(keys ...string) {
	for _, k := range keys {
		_ = os.Unsetenv(k)
	}
}

// TestNewConfigWithDefaults ensures required fields missing cause validation error.
func TestNewConfigWithDefaults(t *testing.T) {
	_, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)
}

// TestNewConfigWithEnvFile verifies config loads correctly from .env file.
func TestNewConfigWithEnvFile(t *testing.T) {
	content := []byte(`
	GRPC_SERVER_URL=127.0.0.1:7000
	DATABASE_DRIVER=postgres
	REDIS_ADDR=localhost:6380
	TRACING_SERVICE_NAME=orders-service
	DATABASE_POOL_MAX_LIFETIME=10m`)

	tmpFile, err := os.CreateTemp("", "test.env")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	tmpFile.Close()

	envLoader := func(path string) error {
		return env.Load(path)
	}

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		EnvPath:   tmpFile.Name(),
		EnvLoader: envLoader,
		Logger:    logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	assert.Equal(t, "127.0.0.1:7000", cfg.GRPCServer.URL)
	assert.Equal(t, "localhost:6380", cfg.Redis.Addr)
	assert.Equal(t, "orders-service", cfg.Tracing.ServiceName)
	assert.Equal(t, 10*time.Minute, cfg.Database.PoolConnMaxLifetime)
}

// TestNewConfigWithValidEnv ensures valid env vars produce a valid config.
func TestNewConfigWithValidEnv(t *testing.T) {
	clearEnv(
		"GRPC_SERVER_URL", "DATABASE_DRIVER", "REDIS_ADDR", "TRACING_SERVICE_NAME",
	)

	os.Setenv("GRPC_SERVER_URL", "localhost:6001")
	os.Setenv("DATABASE_DRIVER", "mysql")
	os.Setenv("REDIS_ADDR", "redis:6379")
	os.Setenv("TRACING_SERVICE_NAME", "user-service")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	assert.Equal(t, "localhost:6001", cfg.GRPCServer.URL)
	assert.Equal(t, "mysql", cfg.Database.Driver)
	assert.Equal(t, "redis:6379", cfg.Redis.Addr)
	assert.Equal(t, "user-service", cfg.Tracing.ServiceName)
}

// TestNewConfigWithInvalidDriver ensures unsupported driver fails validation.
func TestNewConfigWithInvalidDriver(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_DSN", "DATABASE_DRIVER")

	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_DRIVER", "oracle") // invalid

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.Error(t, err)
	require.Nil(t, cfg)
}

// TestNewConfigWithDefaultsApplied ensures defaults are applied for optional fields.
func TestNewConfigWithDefaultsApplied(t *testing.T) {
	clearEnv(
		"GRPC_SERVER_URL", "HTTP_SERVER_URL", "DATABASE_DSN", "DATABASE_DRIVER",
		"REDIS_ADDR", "TRACING_SERVICE_NAME",
	)

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	assert.Equal(t, ":5000", cfg.GRPCServer.URL)
	assert.Equal(t, ":4000", cfg.HTTPServer.URL)
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "go-microservice-boilerplate", cfg.Tracing.ServiceName)
}
