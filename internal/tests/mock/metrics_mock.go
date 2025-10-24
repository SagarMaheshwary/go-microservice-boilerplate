package mock

import (
	"github.com/prometheus/client_golang/prometheus"
	testifymock "github.com/stretchr/testify/mock"
)

type MockCollector struct {
	testifymock.Mock
	Registered bool
}

func (m *MockCollector) Register(r *prometheus.Registry) {
	m.Registered = true
	m.Called(r)
}
