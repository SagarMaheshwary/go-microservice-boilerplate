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
	require.Error(t, err)
}

// TestNewConfigWithEnvFile verifies config loads correctly from .env file.
func TestNewConfigWithEnvFile(t *testing.T) {
	content := []byte(`
	GRPC_SERVER_URL=127.0.0.1:6000
	DATABASE_DSN=postgres://user:pass@localhost:5432/envdb
	DATABASE_DRIVER=postgres
	DATABASE_POOL_MAX_IDLE=5
	DATABASE_POOL_MAX_OPEN=20
	DATABASE_POOL_MAX_LIFETIME=15m
	REDIS_ADDR=localhost:6379
	REDIS_PASSWORD=secret
	REDIS_DB=2
	REDIS_DIAL_TIMEOUT=2s
	REDIS_READ_TIMEOUT=1s
	REDIS_WRITE_TIMEOUT=1s
	REDIS_POOL_SIZE=30
	REDIS_MIN_IDLE_CONNECTIONS=10`)

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

	// Database
	assert.Equal(t, "127.0.0.1:6000", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/envdb", cfg.Database.DSN)
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, 5, cfg.Database.PoolMaxIdleConns)
	assert.Equal(t, 20, cfg.Database.PoolMaxOpenConns)
	assert.Equal(t, 15*time.Minute, cfg.Database.PoolConnMaxLifetime)

	// Redis
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.Equal(t, 2*time.Second, cfg.Redis.DialTimeout)
	assert.Equal(t, 1*time.Second, cfg.Redis.ReadTimeout)
	assert.Equal(t, 1*time.Second, cfg.Redis.WriteTimeout)
	assert.Equal(t, 30, cfg.Redis.PoolSize)
	assert.Equal(t, 10, cfg.Redis.MinIdleConns)
}

// TestNewConfigWithValidEnv ensures valid env vars produce a valid config.
func TestNewConfigWithValidEnv(t *testing.T) {
	clearEnv(
		"GRPC_SERVER_URL", "DATABASE_DSN", "DATABASE_DRIVER",
		"DATABASE_POOL_MAX_IDLE", "DATABASE_POOL_MAX_OPEN", "DATABASE_POOL_MAX_LIFETIME",
		"REDIS_ADDR", "REDIS_PASSWORD", "REDIS_DB", "REDIS_DIAL_TIMEOUT",
		"REDIS_READ_TIMEOUT", "REDIS_WRITE_TIMEOUT", "REDIS_POOL_SIZE", "REDIS_MIN_IDLE_CONNECTIONS",
	)

	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_DRIVER", "mysql")
	os.Setenv("DATABASE_POOL_MAX_IDLE", "3")
	os.Setenv("DATABASE_POOL_MAX_OPEN", "12")
	os.Setenv("DATABASE_POOL_MAX_LIFETIME", "45s")

	os.Setenv("REDIS_ADDR", "localhost:6380")
	os.Setenv("REDIS_PASSWORD", "mypassword")
	os.Setenv("REDIS_DB", "1")
	os.Setenv("REDIS_DIAL_TIMEOUT", "4s")
	os.Setenv("REDIS_READ_TIMEOUT", "2s")
	os.Setenv("REDIS_WRITE_TIMEOUT", "2s")
	os.Setenv("REDIS_POOL_SIZE", "25")
	os.Setenv("REDIS_MIN_IDLE_CONNECTIONS", "7")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	// Database
	assert.Equal(t, "localhost:50051", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Database.DSN)
	assert.Equal(t, "mysql", cfg.Database.Driver)
	assert.Equal(t, 3, cfg.Database.PoolMaxIdleConns)
	assert.Equal(t, 12, cfg.Database.PoolMaxOpenConns)
	assert.Equal(t, 45*time.Second, cfg.Database.PoolConnMaxLifetime)

	// Redis
	assert.Equal(t, "localhost:6380", cfg.Redis.Addr)
	assert.Equal(t, "mypassword", cfg.Redis.Password)
	assert.Equal(t, 1, cfg.Redis.DB)
	assert.Equal(t, 4*time.Second, cfg.Redis.DialTimeout)
	assert.Equal(t, 2*time.Second, cfg.Redis.ReadTimeout)
	assert.Equal(t, 2*time.Second, cfg.Redis.WriteTimeout)
	assert.Equal(t, 25, cfg.Redis.PoolSize)
	assert.Equal(t, 7, cfg.Redis.MinIdleConns)
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
		"GRPC_SERVER_URL", "DATABASE_DSN", "DATABASE_DRIVER",
		"DATABASE_POOL_MAX_IDLE", "DATABASE_POOL_MAX_OPEN", "DATABASE_POOL_MAX_LIFETIME",
		"REDIS_ADDR", "REDIS_PASSWORD", "REDIS_DB", "REDIS_DIAL_TIMEOUT",
		"REDIS_READ_TIMEOUT", "REDIS_WRITE_TIMEOUT", "REDIS_POOL_SIZE", "REDIS_MIN_IDLE_CONNECTIONS",
	)

	// Required only
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	// Defaults
	assert.Equal(t, ":5000", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, 10, cfg.Database.PoolMaxIdleConns)
	assert.Equal(t, 100, cfg.Database.PoolMaxOpenConns)
	assert.Equal(t, time.Hour, cfg.Database.PoolConnMaxLifetime)

	assert.Equal(t, "", cfg.Redis.Addr)
	assert.Equal(t, "default", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, 5*time.Second, cfg.Redis.DialTimeout)
	assert.Equal(t, 3*time.Second, cfg.Redis.ReadTimeout)
	assert.Equal(t, 3*time.Second, cfg.Redis.WriteTimeout)
	assert.Equal(t, 20, cfg.Redis.PoolSize)
	assert.Equal(t, 5, cfg.Redis.MinIdleConns)
}
