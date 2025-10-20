package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Close() error
}

type RedisCache struct {
	client *redis.Client
}

type Opts struct {
	Config *config.Redis
	Logger logger.Logger
}

func NewRedisCache(ctx context.Context, opts *Opts) (CacheService, error) {
	cfg := opts.Config
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	r := &RedisCache{client: rdb}

	if err := r.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	opts.Logger.Info("Redis connected")

	return r, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	if r == nil || r.client == nil {
		return fmt.Errorf("cannot close: redis is not initialized")
	}
	return r.client.Close()
}
