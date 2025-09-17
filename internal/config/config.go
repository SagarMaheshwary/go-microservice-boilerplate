package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

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
	URL      string `validate:"required,url"`
	PoolSize int    `validate:"required,gt=0"`
}

func NewConfig() (*Config, error) {
	return NewConfigWithOptions(LoaderOptions{
		EnvPath: path.Join(rootDir(), "..", ".env"),
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
		log.Info("Loaded environment variables from %q", opts.EnvPath)
	} else {
		log.Info(".env file not found, using system environment variables")
	}

	cfg := &Config{
		GRPCServer: &GRPCServer{
			URL: getEnv("GRPC_SERVER_URL", ":5002"),
		},
		Database: &Database{
			URL:      getEnv("DATABASE_URL", ""),
			PoolSize: getEnvInt("DATABASE_POOL_SIZE", 4),
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
