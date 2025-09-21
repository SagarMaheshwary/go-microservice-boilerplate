package database_test

import (
	"os"
	"testing"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/stretchr/testify/assert"
)

// Test that a supported driver (sqlite) creates a working DB connection
func TestNew_SupportedDriver_SQLite(t *testing.T) {
	log := logger.NewZerologLogger("info", os.Stderr)

	db, err := database.NewDatabase(&database.Opts{
		Config: &config.Database{DSN: ":memory:", Driver: "sqlite"},
		Logger: log,
	})
	assert.NoError(t, err, "expected no error for sqlite driver")
	assert.NotNil(t, db, "expected db instance to be not nil")

	// Check DB() method returns gorm.DB
	gormDB := db.DB()
	assert.NotNil(t, gormDB, "expected gorm.DB instance")

	// Test Close() does not error
	err = db.Close()
	assert.NoError(t, err, "expected no error on close")
}

// Test unsupported driver returns error
func TestNew_UnsupportedDriver(t *testing.T) {
	log := logger.NewZerologLogger("info", os.Stderr)
	db, err := database.NewDatabase(&database.Opts{
		Config: &config.Database{DSN: "some-dsn", Driver: "mongo"},
		Logger: log,
	})
	assert.Error(t, err, "expected error for unsupported driver")
	assert.Nil(t, db, "db should be nil on unsupported driver")
}

// Test invalid DSN
func TestNew_InvalidDSN(t *testing.T) {
	log := logger.NewZerologLogger("info", os.Stderr)
	// For sqlite, an empty string is invalid
	db, err := database.NewDatabase(&database.Opts{
		Config: &config.Database{DSN: "", Driver: "sqlite"},
		Logger: log,
	})
	assert.Error(t, err, "expected error for invalid DSN")
	assert.Nil(t, db, "db should be nil on invalid DSN")
}

// Test Close() on invalid sql.DB (when DB() fails internally)
func TestDatabase_Close_InvalidDB(t *testing.T) {
	db := &database.Database{}
	err := db.Close()
	assert.Error(t, err, "expected error when closing nil DB")
}
