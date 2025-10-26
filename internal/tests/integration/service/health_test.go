package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/tests/testutils"
)

func TestHealthService_Check_Success(t *testing.T) {
	db := testutils.SetupPostgres(t)
	redis := testutils.SetupRedis(t)

	svc := service.NewHealthService(&service.HealthServiceOpts{
		Database: db,
		Cache:    redis,
	})

	status := svc.Check(context.Background())

	require.NotNil(t, status)
	assert.Equal(t, "ready", status.Status)
	assert.Equal(t, "ok", status.Details["database"])
	assert.Equal(t, "ok", status.Details["cache"])
}

func TestHealthService_Check_DatabaseDown(t *testing.T) {
	db := testutils.SetupPostgres(t)
	redis := testutils.SetupRedis(t)

	// Simulate broken DB connection
	sqlDB, err := db.DB().DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close()) // forcibly close connection

	svc := service.NewHealthService(&service.HealthServiceOpts{
		Database: db,
		Cache:    redis,
	})

	status := svc.Check(context.Background())

	require.NotNil(t, status)
	assert.Equal(t, "unready", status.Status)
	assert.Equal(t, "ok", status.Details["cache"])
	assert.NotEqual(t, "ok", status.Details["database"])
}

func TestHealthService_Check_CacheDown(t *testing.T) {
	db := testutils.SetupPostgres(t)
	redis := testutils.SetupRedis(t)

	// Simulate broken cache connection by closing redis client
	require.NoError(t, redis.Close())

	svc := service.NewHealthService(&service.HealthServiceOpts{
		Database: db,
		Cache:    redis,
	})

	status := svc.Check(context.Background())

	require.NotNil(t, status)
	assert.Equal(t, "unready", status.Status)
	assert.Equal(t, "ok", status.Details["database"])
	assert.NotEqual(t, "ok", status.Details["cache"])
}
