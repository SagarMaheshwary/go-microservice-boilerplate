package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/model"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/tests/testutils"
)

func TestUserService_FindByID(t *testing.T) {
	db := testutils.SetupPostgres(t)

	// Seed test data
	u := &model.User{Name: "Alice", Email: "alice@example.com"}
	require.NoError(t, db.DB().Create(u).Error)

	userService := service.NewUserService(db)

	got, err := userService.FindByID(context.Background(), u.ID)
	require.NoError(t, err)

	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
}
