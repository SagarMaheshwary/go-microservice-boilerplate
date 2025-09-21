package database

import (
	"fmt"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseService interface {
	DB() *gorm.DB
	Close() error
}

type Database struct {
	db     *gorm.DB
	Logger logger.Logger
}

type Opts struct {
	Config *config.Database
	Logger logger.Logger
}

func NewDatabase(opts *Opts) (DatabaseService, error) {
	var (
		db  *gorm.DB
		err error
	)

	driver := opts.Config.Driver
	dsn := opts.Config.DSN

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		//gorm will create a db connection even if dsn is empty so adding this check to
		//keep the connection flow consistent
		if opts.Config.DSN == "" {
			return nil, fmt.Errorf("invalid DSN: sqlite requires a non-empty DSN")
		}

		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver %s", driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", driver, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(opts.Config.PoolMaxIdleConns)
	sqlDB.SetMaxOpenConns(opts.Config.PoolMaxOpenConns)
	sqlDB.SetConnMaxLifetime(opts.Config.PoolConnMaxLifetime)

	opts.Logger.Info(`Database connected on driver "%s"`, driver)

	return &Database{db: db, Logger: opts.Logger}, nil
}

func (d *Database) DB() *gorm.DB {
	return d.db
}

func (d *Database) Close() error {
	if d == nil || d.db == nil {
		return fmt.Errorf("cannot close: database is not initialized")
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
