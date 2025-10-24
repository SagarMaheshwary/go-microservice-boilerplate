package interceptor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/metrics"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server/interceptor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupTestRegistry(t *testing.T) *prometheus.Registry {
	t.Helper()
	reg := prometheus.NewRegistry()
	reg.MustRegister(metrics.GRPCRequestCounter, metrics.GRPCRequestLatency)
	metrics.GRPCRequestCounter.Reset()
	metrics.GRPCRequestLatency.Reset()
	return reg
}

func assertCounterAndLatency(
	t *testing.T,
	reg *prometheus.Registry,
	method string,
	expectedCode string,
) {
	t.Helper()
	mfs, err := reg.Gather()
	require.NoError(t, err)

	var counterFound, latencyFound bool

	for _, mf := range mfs {
		switch mf.GetName() {
		case "grpc_requests_total":
			for _, metric := range mf.GetMetric() {
				methodLabel := metric.GetLabel()[0].GetValue()
				statusLabel := metric.GetLabel()[1].GetValue()
				if methodLabel == method && statusLabel == expectedCode {
					counterFound = true
				}
			}
		case "grpc_request_duration_seconds":
			for _, metric := range mf.GetMetric() {
				methodLabel := metric.GetLabel()[0].GetValue()
				if methodLabel == method {
					h := metric.GetHistogram()
					if h.GetSampleCount() >= 1 {
						latencyFound = true
					}
				}
			}
		}
	}

	assert.True(t, counterFound, "grpc_requests_total counter should be incremented with correct labels")
	assert.True(t, latencyFound, "grpc_request_duration_seconds histogram should be observed")
}

func TestMetricsInterceptor_Success(t *testing.T) {
	reg := setupTestRegistry(t)
	interceptorFn := interceptor.MetricsInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(5 * time.Millisecond)
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Success"}
	resp, err := interceptorFn(context.Background(), "req", info, handler)

	require.NoError(t, err)
	assert.Equal(t, "ok", resp)

	assertCounterAndLatency(t, reg, info.FullMethod, "OK")
}

func TestMetricsInterceptor_InternalError(t *testing.T) {
	reg := setupTestRegistry(t)
	interceptorFn := interceptor.MetricsInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Error(codes.Internal, "oops")
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Internal"}
	resp, err := interceptorFn(context.Background(), "req", info, handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	assertCounterAndLatency(t, reg, info.FullMethod, "Internal")
}

func TestMetricsInterceptor_CustomError(t *testing.T) {
	reg := setupTestRegistry(t)
	interceptorFn := interceptor.MetricsInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("custom")
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Custom"}
	resp, err := interceptorFn(context.Background(), "req", info, handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	assertCounterAndLatency(t, reg, info.FullMethod, "Unknown")
}
