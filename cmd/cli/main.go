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

	cfg, err := config.NewConfig(log)
	if err != nil {
		log.Fatal(err.Error())
	}

	switch cmd {
	case "seed":
		if len(os.Args) > 2 {
			cfg.Database.DSN = os.Args[2]
			log.Info("Using DSN from argument")
		}

		db, err := database.NewDatabase(&database.Opts{
			Config: cfg.Database,
			Logger: log,
		})
		if err != nil {
			log.Fatal(err.Error())
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Error(err.Error())
			}
		}()

		err = seeder.RunAll(&seeder.Opts{
			DB:  db.DB(),
			Log: log,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	default:
		log.Error("Unknown command " + cmd)
	}
}
