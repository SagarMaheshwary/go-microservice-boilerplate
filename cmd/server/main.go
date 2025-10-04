package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server"
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

	grpcServer := server.NewServer(&server.Opts{
		Config:   cfg.GRPCServer,
		Logger:   log,
		Database: db,
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

	log.Info("Shutdown complete!")
}
