package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

type Config struct {
	GRPCServer *GRPCServer
	Database   *Database
}

type GRPCServer struct {
	URL string
}

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
}

type LoaderOptions struct {
	EnvPath     string
	EnvLoader   func(string) error
	FileChecker func(string) bool
	Logger      logger.Logger
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
		envLoader = func(path string) error { return env.Load(path) }
	}
	fileChecker := opts.FileChecker
	if fileChecker == nil {
		fileChecker = func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		}
	}

	if opts.EnvPath != "" && fileChecker(opts.EnvPath) {
		if err := envLoader(opts.EnvPath); err != nil {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}
		log.Info("Loaded environment variables from %q", opts.EnvPath)
	} else {
		log.Info(".env file not found, using system environment variables")
	}

	return &Config{
		GRPCServer: &GRPCServer{
			URL: getEnv("GRPC_SERVER_URL", ":5002"),
		},
		Database: &Database{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvInt("DATABASE_PORT", 5432),
			Username: getEnv("DATABASE_USERNAME", "postgres"),
			Password: getEnv("DATABASE_PASSWORD", "password"),
		},
	}, nil
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
