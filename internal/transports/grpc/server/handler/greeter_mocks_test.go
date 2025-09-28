package handler_test

import (
	"context"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/model"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	var user *model.User
	if u := args.Get(0); u != nil {
		user = u.(*model.User)
	}
	return user, args.Error(1)
}
