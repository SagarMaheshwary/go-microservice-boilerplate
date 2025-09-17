package main

import (
	"os"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

func main() {
	log := logger.NewZerologLogger("info", os.Stderr)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("%v", err)
	}

	log.Info("Hello World")
	log.Info("GRPC server running on %s", cfg.GRPCServer.URL)
}
