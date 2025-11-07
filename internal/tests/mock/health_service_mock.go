package mock

import (
	"context"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) Check(ctx context.Context) service.HealthStatus {
	args := m.Called(ctx)
	return args.Get(0).(service.HealthStatus)
}

func (m *MockHealthService) SetReady(isReady bool) {
	m.Called(isReady)
}
