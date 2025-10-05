package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

type LoaderOptions struct {
	EnvPath   string
	EnvLoader func(string) error
	Logger    logger.Logger
}

type Config struct {
	GRPCServer *GRPCServer `validate:"required"`
	Database   *Database   `validate:"required"`
}

type GRPCServer struct {
	URL string `validate:"required,hostname_port"`
}

type Database struct {
	DSN                 string        `validate:"required"`
	Driver              string        `validate:"required,oneof=postgres mysql sqlite"`
	PoolMaxIdleConns    int           `validate:"gte=0"`
	PoolMaxOpenConns    int           `validate:"gte=0"`
	PoolConnMaxLifetime time.Duration `validate:"gte=0"` // must be non-negative
}

func NewConfig(log logger.Logger) (*Config, error) {
	return NewConfigWithOptions(LoaderOptions{
		EnvPath: path.Join(rootDir(), "..", ".env"),
		Logger:  log,
	})
}

func NewConfigWithOptions(opts LoaderOptions) (*Config, error) {
	log := opts.Logger
	if log == nil {
		log = logger.NewZerologLogger("info", os.Stderr)
	}

	envLoader := opts.EnvLoader
	if envLoader == nil {
		envLoader = func(path string) error {
			_, err := os.Stat(path)
			if err != nil {
				return err
			}

			return env.Load(path)
		}
	}

	if err := envLoader(opts.EnvPath); err == nil {
		log.Info("Loaded environment variables from" + opts.EnvPath)
	} else {
		log.Info("failed to load .env file, using system environment variables")
	}

	cfg := &Config{
		GRPCServer: &GRPCServer{
			URL: getEnv("GRPC_SERVER_URL", ":5000"),
		},
		Database: &Database{
			DSN:                 getEnv("DATABASE_DSN", ""),
			Driver:              getEnv("DATABASE_DRIVER", "postgres"),
			PoolMaxIdleConns:    getEnvInt("DATABASE_POOL_MAX_IDLE", 10),
			PoolMaxOpenConns:    getEnvInt("DATABASE_POOL_MAX_OPEN", 100),
			PoolConnMaxLifetime: getEnvDuration("DATABASE_POOL_MAX_LIFETIME", time.Hour),
		},
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d)
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}
