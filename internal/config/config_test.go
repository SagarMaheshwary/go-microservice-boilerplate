package config_test

import (
	"os"
	"testing"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// clearEnv is a helper that removes the provided environment variables.
// This ensures each test runs in a clean environment.
func clearEnv(keys ...string) {
	for _, k := range keys {
		os.Unsetenv(k)
	}
}

// TestNewConfigWithDefaults verifies that if required environment variables
// are missing, config validation fails.
func TestNewConfigWithDefaults(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_URL", "DATABASE_POOL_SIZE")

	_, err := config.NewConfigWithOptions(config.LoaderOptions{})
	require.Error(t, err)
}

// TestNewConfigWithEnvFile verifies that configuration can be loaded from
// a .env file using a custom EnvLoader.
func TestNewConfigWithEnvFile(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_URL", "DATABASE_POOL_SIZE")

	// Create a temporary .env file with test values.
	content := []byte(`
	GRPC_SERVER_URL=127.0.0.1:6000
	DATABASE_URL=postgres://user:pass@localhost:5432/envdb
	DATABASE_POOL_SIZE=12`)

	tmpFile, err := os.CreateTemp("", "test.env")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	tmpFile.Close()

	// Use the real env.Load for this test.
	envLoader := func(path string) error {
		return env.Load(path)
	}

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{
		EnvPath:   tmpFile.Name(),
		EnvLoader: envLoader,
	})
	require.NoError(t, err)

	assert.Equal(t, "127.0.0.1:6000", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/envdb", cfg.Database.URL)
	assert.Equal(t, 12, cfg.Database.PoolSize)
}

// TestNewConfigWithValidEnv verifies that configuration loads successfully
// when valid environment variables are present.
func TestNewConfigWithValidEnv(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_URL", "DATABASE_POOL_SIZE")

	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_POOL_SIZE", "15")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{})
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "localhost:50051", cfg.GRPCServer.URL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Database.URL)
	assert.Equal(t, 15, cfg.Database.PoolSize)
}

// TestNewConfigWithInvalidDatabaseURL verifies that validation fails if
// DATABASE_URL is not a valid URL.
func TestNewConfigWithInvalidDatabaseURL(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_URL", "DATABASE_POOL_SIZE")

	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_URL", "://not-a-url") // invalid
	os.Setenv("DATABASE_POOL_SIZE", "10")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{})
	require.Error(t, err)
	require.Nil(t, cfg)
}

// TestNewConfigWithInvalidPoolSize verifies that validation fails if
// DATABASE_POOL_SIZE does not meet the "gt=0" constraint.
func TestNewConfigWithInvalidPoolSize(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_URL", "DATABASE_POOL_SIZE")

	os.Setenv("GRPC_SERVER_URL", "localhost:50051")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DATABASE_POOL_SIZE", "0") // invalid (must be > 0)

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{})
	require.Error(t, err)
	require.Nil(t, cfg)
}

// TestNewConfigWithDefaultsApplied verifies that default values are applied
// when optional environment variables are missing.
func TestNewConfigWithDefaultsApplied(t *testing.T) {
	clearEnv("GRPC_SERVER_URL", "DATABASE_URL", "DATABASE_POOL_SIZE")

	// Only set the required DATABASE_URL.
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg, err := config.NewConfigWithOptions(config.LoaderOptions{})
	require.NoError(t, err)

	// GRPC_SERVER_URL should fallback to default ":5002".
	assert.Equal(t, ":5002", cfg.GRPCServer.URL)
	// DATABASE_POOL_SIZE should fallback to 4.
	assert.Equal(t, 4, cfg.Database.PoolSize)
}
