package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server"
	helloworld "github.com/sagarmaheshwary/go-microservice-boilerplate/proto/hello_world"
)

// TestSayHello_Success verifies that the SayHello RPC
// returns the expected response without error.
func TestSayHello_Success(t *testing.T) {
	s := &server.GreeterServer{}

	req := &helloworld.SayHelloRequest{}

	resp, err := s.SayHello(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Hello, World!", resp.Message)
}

// TestSayHello_EmptyName verifies that even if the request has no name,
// the RPC still returns a valid response (current implementation ignores Name).
func TestSayHello_EmptyName(t *testing.T) {
	s := &server.GreeterServer{}

	req := &helloworld.SayHelloRequest{} // no name

	resp, err := s.SayHello(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Hello, World!", resp.Message)
}
