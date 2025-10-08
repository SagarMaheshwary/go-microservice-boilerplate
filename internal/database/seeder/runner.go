package seeder

import (
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"gorm.io/gorm"
)

type SeederFunc struct {
	Name string
	Func func(db *gorm.DB) error
}

var defaultSeeders = []SeederFunc{
	{Name: "SeedUsers", Func: SeedUsers},
	// Add more seeders here
}

type Opts struct {
	DB      *gorm.DB
	Log     logger.Logger
	Seeders []SeederFunc
}

func RunAll(opts *Opts) error {
	log := opts.Log
	seeders := opts.Seeders
	if len(seeders) == 0 {
		seeders = defaultSeeders
	}

	for _, s := range seeders {
		log.Info("[Seeder] Running " + s.Name)

		if err := s.Func(opts.DB); err != nil {
			log.Info("[Seeder] Failed "+s.Name, logger.Field{Key: "error", Value: err.Error()})
			return err
		}

		log.Info("[Seeder] Completed " + s.Name)
	}

	log.Info("[Seeder] All seeders completed successfully")
	return nil
}
