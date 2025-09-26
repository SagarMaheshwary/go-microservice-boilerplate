package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"google.golang.org/grpc"
)

func LoggerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		elapsedMs := fmt.Sprintf("%.2fms", time.Since(start).Seconds()*1000)

		if err == nil {
			log.Info("gRPC request completed",
				logger.Field{Key: "method", Value: info.FullMethod},
				logger.Field{Key: "duration", Value: elapsedMs},
			)
		} else {
			log.Error("gRPC request failed",
				logger.Field{Key: "method", Value: info.FullMethod},
				logger.Field{Key: "duration", Value: elapsedMs},
				logger.Field{Key: "error", Value: err.Error()},
			)
		}

		return resp, err
	}
}
