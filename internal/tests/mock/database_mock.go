package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) DB() *gorm.DB {
	args := m.Called()
	if db := args.Get(0); db != nil {
		return db.(*gorm.DB)
	}
	return nil
}

func (m *MockDatabase) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}
