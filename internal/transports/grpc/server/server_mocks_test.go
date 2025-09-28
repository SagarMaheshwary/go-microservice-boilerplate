package server_test

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDatabaseService struct {
	mock.Mock
}

func (m *MockDatabaseService) DB() *gorm.DB {
	args := m.Called()
	if db := args.Get(0); db != nil {
		return db.(*gorm.DB)
	}
	return nil
}

func (m *MockDatabaseService) Close() error {
	args := m.Called()
	return args.Error(0)
}
