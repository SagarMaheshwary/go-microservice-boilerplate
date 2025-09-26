package interceptor_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server/interceptor"
)

// helper to parse zerolog JSON into a map
func parseLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	var entry map[string]any
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	return entry
}

// TestLoggingInterceptor_Success verifies that when the gRPC handler
// completes successfully, the interceptor logs at "info" level with
// the method name and request duration.
func TestLoggingInterceptor_Success(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	// interceptor under test
	interceptorFn := interceptor.LoggerInterceptor(log)

	// fake handler returns success
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// invoke interceptor
	resp, err := interceptorFn(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		handler,
	)

	// assertions
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)

	entry := parseLog(t, &buf)
	assert.Equal(t, "info", entry["level"])
	assert.Equal(t, "gRPC request completed", entry["message"])
	assert.Equal(t, "/test.Service/Method", entry["method"])
	assert.Contains(t, entry["duration"], "ms")
}

// TestLoggingInterceptor_Error verifies that when the gRPC handler
// returns an error, the interceptor logs at "error" level with
// the method name, request duration, and the error message.
func TestLoggingInterceptor_Error(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	interceptorFn := interceptor.LoggerInterceptor(log)

	// fake handler returns error
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("boom")
	}

	// invoke interceptor
	resp, err := interceptorFn(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		handler,
	)

	// assertions
	require.Error(t, err)
	assert.Nil(t, resp)

	entry := parseLog(t, &buf)
	assert.Equal(t, "error", entry["level"])
	assert.Equal(t, "gRPC request failed", entry["message"])
	assert.Equal(t, "/test.Service/Method", entry["method"])
	assert.Contains(t, entry["duration"], "ms")
	assert.Equal(t, "boom", entry["error"])
}
