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

// clearEnv is a helper that removes the provided environment variables.
func clearEnv(keys ...string) {
	log := logger.NewZerologLogger("info", io.Discard)
	for _, k := range keys {
		if err := os.Unsetenv(k); err != nil {
			log.Error(err.Error())
		}
	}
}

// TestNewConfigWithDefaults ensures required fields missing cause validation error.
func TestNewConfigWithDefaults(t *testing.T) {
	_, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.Error(t, err)
}

// TestNewConfigWithEnvFile verifies config loads from .env file.
func TestNewConfigWithEnvFile(t *testing.T) {
	// Create a temporary .env file
	content := []byte(`
	GRPC_SERVER_URL=127.0.0.1:6000
	DATABASE_DSN=postgres://user:pass@localhost:5432/envdb
	DATABASE_DRIVER=postgres
	DATABASE_POOL_MAX_IDLE=5
	DATABASE_POOL_MAX_OPEN=20
	DATABASE_POOL_MAX_LIFETIME=15m`)

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

	assert.Equal(t, "127.0.0.1:6000", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/envdb", cfg.Database.DSN)
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, 5, cfg.Database.PoolMaxIdleConns)
	assert.Equal(t, 20, cfg.Database.PoolMaxOpenConns)
	assert.Equal(t, 15*time.Minute, cfg.Database.PoolConnMaxLifetime)
}

// TestNewConfigWithValidEnv ensures valid env vars produce a valid config.
func TestNewConfigWithValidEnv(t *testing.T) {
	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_DRIVER", "mysql")
	os.Setenv("DATABASE_POOL_MAX_IDLE", "3")
	os.Setenv("DATABASE_POOL_MAX_OPEN", "12")
	os.Setenv("DATABASE_POOL_MAX_LIFETIME", "45s")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	assert.Equal(t, "localhost:50051", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Database.DSN)
	assert.Equal(t, "mysql", cfg.Database.Driver)
	assert.Equal(t, 3, cfg.Database.PoolMaxIdleConns)
	assert.Equal(t, 12, cfg.Database.PoolMaxOpenConns)
	assert.Equal(t, 45*time.Second, cfg.Database.PoolConnMaxLifetime)
}

// TestNewConfigWithInvalidDriver ensures unsupported driver fails validation.
func TestNewConfigWithInvalidDriver(t *testing.T) {
	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_DRIVER", "oracle") // not allowed

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.Error(t, err)
	require.Nil(t, cfg)
}

// TestNewConfigWithInvalidPoolLifetime ensures invalid duration falls back.
func TestNewConfigWithInvalidPoolLifetime(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_DSN", "DATABASE_DRIVER", "DATABASE_POOL_MAX_IDLE", "DATABASE_POOL_MAX_OPEN", "DATABASE_POOL_MAX_LIFETIME")

	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_POOL_MAX_LIFETIME", "notaduration")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	// Should fall back to default (1h from your config code)
	assert.Equal(t, time.Hour, cfg.Database.PoolConnMaxLifetime)
}

// TestNewConfigWithDefaultsApplied ensures defaults are applied for optional fields.
func TestNewConfigWithDefaultsApplied(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_DSN", "DATABASE_DRIVER", "DATABASE_POOL_MAX_IDLE", "DATABASE_POOL_MAX_OPEN", "DATABASE_POOL_MAX_LIFETIME")

	// Only set required DATABASE_DSN.
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	assert.Equal(t, ":5002", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres", cfg.Database.Driver) // default
	assert.Equal(t, 10, cfg.Database.PoolMaxIdleConns)
	assert.Equal(t, 100, cfg.Database.PoolMaxOpenConns)
	assert.Equal(t, time.Hour, cfg.Database.PoolConnMaxLifetime)
}
