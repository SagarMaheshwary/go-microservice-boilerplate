package seeder_test

import (
	"errors"
	"io"
	"testing"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/seeder"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	log := logger.NewZerologLogger("info", io.Discard)
	db, err := database.NewDatabase(&database.Opts{
		Config: &config.Database{DSN: ":memory:", Driver: "sqlite"},
		Logger: log,
	})

	assert.NoError(t, err)
	return db.DB()
}

// TestRunAll_Success ensures that RunAll executes all provided seeders successfully
// and returns no error when all seeder functions complete without failure.
func TestRunAll_Success(t *testing.T) {
	db := newTestDB(t)
	log := logger.NewZerologLogger("info", io.Discard)

	seeders := []seeder.SeederFunc{
		{Name: "TestSeeder", Func: func(db *gorm.DB) error {
			return nil
		}},
	}

	err := seeder.RunAll(&seeder.Opts{
		DB:      db,
		Log:     log,
		Seeders: seeders,
	})

	assert.NoError(t, err, "expected no error when all seeders succeed")
}

// TestRunAll_Failure verifies that RunAll stops execution and returns an error
// when any seeder function fails. It ensures proper error propagation.
func TestRunAll_Failure(t *testing.T) {
	db := newTestDB(t)
	log := logger.NewZerologLogger("info", io.Discard)

	seeders := []seeder.SeederFunc{
		{Name: "FailingSeeder", Func: func(db *gorm.DB) error {
			return errors.New("boom")
		}},
	}

	err := seeder.RunAll(&seeder.Opts{
		DB:      db,
		Log:     log,
		Seeders: seeders,
	})

	assert.Error(t, err, "expected error when a seeder fails")
	assert.EqualError(t, err, "boom", "expected specific error message to propagate")
}
