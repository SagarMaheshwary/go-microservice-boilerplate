package testutils

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

func SetupPostgres(t *testing.T) database.DatabaseService {
	ctx := context.Background()
	log := logger.NewZerologLogger("info", io.Discard)

	dbName := "testdb"
	username := "test"
	password := "test"

	pgContainer, err := pgcontainer.Run(ctx,
		"postgres:15-alpine",
		pgcontainer.WithDatabase(dbName),
		pgcontainer.WithUsername(username),
		pgcontainer.WithPassword(password),
		pgcontainer.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s ", err.Error())
		}
	})

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port.Port(),
		dbName,
	)

	db, err := database.NewDatabase(&database.Opts{
		Config: &config.Database{
			DSN:                 dsn,
			Driver:              "postgres",
			PoolMaxIdleConns:    10,
			PoolMaxOpenConns:    100,
			PoolConnMaxLifetime: time.Hour,
		},
		Logger: log,
	})
	require.NoError(t, err)

	sqlDB, err := db.DB().DB()
	require.NoError(t, err)

	// Migration driver
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	require.NoError(t, err)

	// Use embedded migrations
	d, err := iofs.New(database.MigrationsFS, "migrations")
	require.NoError(t, err)

	m, err := migrate.NewWithInstance("iofs", d, dbName, driver)
	require.NoError(t, err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}
