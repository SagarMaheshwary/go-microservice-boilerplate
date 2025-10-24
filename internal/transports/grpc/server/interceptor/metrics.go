package interceptor

import (
	"context"
	"time"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		res, err := handler(ctx, req)

		method := info.FullMethod
		statusCode := status.Code(err).String()
		metrics.GRPCRequestCounter.WithLabelValues(method, statusCode).Inc()
		metrics.GRPCRequestLatency.WithLabelValues(method).Observe(time.Since(start).Seconds())

		return res, err
	}
}
