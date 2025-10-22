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

	httpServer := httpserver.NewServer(&httpserver.Opts{
		Config:   cfg.HTTPServer,
		Logger:   log,
		Database: db,
		Cache:    redisCache,
	})
	go func() {
		err = httpServer.Serve()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			stop()
		}
	}()

	grpcServer := grpcserver.NewServer(&grpcserver.Opts{
		Config:   cfg.GRPCServer,
		Logger:   log,
		Database: db,
		Cache:    redisCache,
	})
	go func() {
		err = grpcServer.Serve()
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			stop()
		}
	}()

	<-ctx.Done()

	log.Warn("Shutdown signal received, closing services!")

	grpcServer.Server.GracefulStop()

	if err := db.Close(); err != nil {
		log.Error("failed to close database client", logger.Field{Key: "error", Value: err.Error()})
	}

	if err := redisCache.Close(); err != nil {
		log.Error("failed to close cache client", logger.Field{Key: "error", Value: err.Error()})
	}

	if err := httpServer.Server.Shutdown(ctx); err != nil {
		log.Error("failed to close http server", logger.Field{Key: "error", Value: err.Error()})
	}

	log.Info("Shutdown complete!")
}
