package main

import (
	"os"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/seeder"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

func main() {
	log := logger.NewZerologLogger("info", os.Stderr)

	if len(os.Args) < 2 {
		log.Info("Usage: go run cmd/cli/main.go seed")
		os.Exit(1)
	}

	cmd := os.Args[1]

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	switch cmd {
	case "seed":
		db, err := database.NewDatabase(&database.Opts{
			Config: cfg.Database,
			Logger: log,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		defer db.Close()

		err = seeder.RunAll(&seeder.Opts{
			DB:  db.DB(),
			Log: log,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	default:
		log.Info("Unknown command " + cmd)
	}
}
