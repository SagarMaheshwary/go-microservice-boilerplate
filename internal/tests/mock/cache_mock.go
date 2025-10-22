package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)

	if err := args.Error(1); err != nil {
		return "", err
	}

	return args.Get(0).(string), nil
}

func (m *MockRedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil

}

func (m *MockRedisCache) Delete(ctx context.Context, keys string) error {
	args := m.Called(ctx, keys)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

func (m *MockRedisCache) Ping(ctx context.Context) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

func (m *MockRedisCache) Close() error {
	return nil
}
