package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/metrics"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/tracing"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	grpcserver "github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server"
	httpserver "github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/http/server"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log := logger.NewZerologLogger("info", os.Stderr)

	cfg, err := config.NewConfig(log)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.NewDatabase(&database.Opts{
		Config: cfg.Database,
		Logger: log,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	redisCache, err := cache.NewRedisCache(ctx, &cache.Opts{
		Config: cfg.Redis,
		Logger: log,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	metricsService := metrics.NewMetricsService(cfg.Metrics)
	healthService := service.NewHealthService(map[string]service.DependencyHealthCheck{
		"database": func(ctx context.Context) error {
			return db.DB().Exec("SELECT 1").Error
		},
		"cache": func(ctx context.Context) error {
			return redisCache.Ping(ctx)
		},
	})

	httpServer := httpserver.NewServer(&httpserver.Opts{
		Config:   cfg.HTTPServer,
		Logger:   log,
		Database: db,
		Cache:    redisCache,
		Metrics:  metricsService,
		Health:   healthService,
	})
	go func() {
		err = httpServer.Serve()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			stop()
		}
	}()

	tracerService, err := tracing.NewTracerService(ctx, &tracing.Opts{
		Config: cfg.Tracing,
		Logger: log,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcServer := grpcserver.NewServer(&grpcserver.Opts{
		Config:   cfg.GRPCServer,
		Logger:   log,
		Database: db,
		Cache:    redisCache,
	})
	go func() {
		err = grpcServer.Serve()
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			healthService.SetReady(false)
			stop()
		}
	}()

	<-ctx.Done()

	log.Warn("Shutdown signal received, closing services!")

	// Mark service as not ready
	healthService.SetReady(false)

	grpcServer.Server.GracefulStop()

	if err := db.Close(); err != nil {
		log.Error("failed to close database client", logger.Field{Key: "error", Value: err.Error()})
	}

	if err := redisCache.Close(); err != nil {
		log.Error("failed to close cache client", logger.Field{Key: "error", Value: err.Error()})
	}

	tracerCtx, tracerCancel := context.WithTimeout(context.Background(), cfg.Tracing.ShutdownTimeout)
	if err := tracerService.Shutdown(tracerCtx); err != nil {
		log.Error("failed to close tracing client", logger.Field{Key: "error", Value: err.Error()})
	}
	tracerCancel()

	// Shut down the health server last so it can continue responding to liveness checks
	// (e.g., /livez) while marking the service as not ready (/readyz) during shutdown.
	httpCtx, httpCancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	if err := httpServer.Server.Shutdown(httpCtx); err != nil {
		log.Error("failed to close http server", logger.Field{Key: "error", Value: err.Error()})
	}
	httpCancel()

	log.Info("Shutdown complete!")
}
