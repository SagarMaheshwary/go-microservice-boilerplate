package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/model"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server/handler"
	helloworld "github.com/sagarmaheshwary/go-microservice-boilerplate/proto/hello_world"
)

// TestSayHello_Success verifies SayHello returns expected response and user.
func TestSayHello_Success(t *testing.T) {
	mockUser := &model.User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	mockService := new(MockUserService)
	mockService.On("FindByID", mock.Anything, uint(1)).
		Return(mockUser, nil)

	s := handler.NewGreeterServer(mockService)

	req := &helloworld.SayHelloRequest{}
	resp, err := s.SayHello(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Hello, World!", resp.Message)
	assert.Equal(t, mockUser.Name, resp.User.Name)
	assert.Equal(t, mockUser.Email, resp.User.Email)

	// verify expectations
	mockService.AssertExpectations(t)
}
