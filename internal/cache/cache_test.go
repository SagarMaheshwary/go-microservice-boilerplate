package cache_test

import (
	"context"
	"io"
	"testing"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/stretchr/testify/assert"
)

// Test that invalid Redis address returns an error
func TestNewRedisCache_InvalidConfig(t *testing.T) {
	ctx := context.Background()
	log := logger.NewZerologLogger("info", io.Discard)

	r, err := cache.NewRedisCache(ctx, &cache.Opts{
		Config: &config.Redis{
			Addr: "invalid:6379", // invalid host
			DB:   0,
		},
		Logger: log,
	})

	assert.Error(t, err, "expected error for invalid Redis address")
	assert.Nil(t, r, "expected nil cache on error")
}

// Test Close() on nil client
func TestRedisCache_Close_InvalidClient(t *testing.T) {
	r := &cache.RedisCache{}
	err := r.Close()
	assert.Error(t, err, "expected error when closing nil client")
}
