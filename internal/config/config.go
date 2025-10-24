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
	HTTPServer *HTTPServer `validate:"required"`
	Database   *Database   `validate:"required"`
	Redis      *Redis      `validate:"required"`
	Metrics    *Metrics    `validate:"required"`
	Tracing    *Tracing    `validate:"required"`
}

type GRPCServer struct {
	URL string `validate:"required,hostname_port"`
}

type HTTPServer struct {
	URL string `validate:"required,hostname_port"`
}

type Database struct {
	DSN                 string        `validate:"required"`
	Driver              string        `validate:"required,oneof=postgres mysql sqlite"`
	PoolMaxIdleConns    int           `validate:"gte=0"`
	PoolMaxOpenConns    int           `validate:"gte=0"`
	PoolConnMaxLifetime time.Duration `validate:"gte=0"` // must be non-negative
}

type Redis struct {
	Addr         string        `validate:"required"`
	Password     string        `validate:"required"`
	DB           int           `validate:"gte=0"`
	DialTimeout  time.Duration `validate:"gte=0"`
	ReadTimeout  time.Duration `validate:"gte=0"`
	WriteTimeout time.Duration `validate:"gte=0"`
	PoolSize     int           `validate:"gte=0"`
	MinIdleConns int           `validate:"gte=0"`
}

type Metrics struct {
	EnableDefaultMetrics bool
}

type Tracing struct {
	ServiceName  string `validate:"required"`
	CollectorURL string `validate:"required,hostname_port"`
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
		HTTPServer: &HTTPServer{
			URL: getEnv("HTTP_SERVER_URL", ":4000"),
		},
		Database: &Database{
			DSN:                 getEnv("DATABASE_DSN", "postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable"),
			Driver:              getEnv("DATABASE_DRIVER", "postgres"),
			PoolMaxIdleConns:    getEnvInt("DATABASE_POOL_MAX_IDLE", 10),
			PoolMaxOpenConns:    getEnvInt("DATABASE_POOL_MAX_OPEN", 100),
			PoolConnMaxLifetime: getEnvDuration("DATABASE_POOL_MAX_LIFETIME", time.Hour),
		},
		Redis: &Redis{
			Addr:         getEnv("REDIS_ADDR", "localhost:6379"),
			Password:     getEnv("REDIS_PASSWORD", "default"),
			DB:           getEnvInt("REDIS_DB", 0),
			DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", time.Second*5),
			ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", time.Second*3),
			WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", time.Second*3),
			PoolSize:     getEnvInt("REDIS_POOL_SIZE", 20),
			MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNECTIONS", 5),
		},
		Metrics: &Metrics{
			EnableDefaultMetrics: getEnvBool("METRICS_ENABLE_DEFAULT_METRICS", false),
		},
		Tracing: &Tracing{
			ServiceName:  getEnv("TRACING_SERVICE_NAME", "go-microservice-boilerplate"),
			CollectorURL: getEnv("TRACING_COLLECTOR_URL", "localhost:4318"),
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

func getEnvBool(key string, defaultVal bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}
