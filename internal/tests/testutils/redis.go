package testutils

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func SetupRedis(t *testing.T) cache.CacheService {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx,
		"redis:7.2",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(redisContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	})

	host, _ := redisContainer.Host(ctx)
	port, _ := redisContainer.MappedPort(ctx, "6379")

	redisCache, err := cache.NewRedisCache(ctx, &cache.Opts{
		Config: &config.Redis{
			Addr:         fmt.Sprintf("%s:%s", host, port.Port()),
			Password:     "",
			DB:           0,
			DialTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
			PoolSize:     20,
			MinIdleConns: 5,
		},
		Logger: logger.NewZerologLogger("info", io.Discard),
	})
	require.NoError(t, err)

	return redisCache
}
