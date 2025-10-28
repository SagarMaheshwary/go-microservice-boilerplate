package server

import (
	"net"
	"net/http"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/observability/metrics"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/http/server/handler"
)

type Opts struct {
	Config   *config.HTTPServer
	Logger   logger.Logger
	Database database.DatabaseService
	Cache    cache.CacheService
	Metrics  metrics.MetricsService
}

type HTTPServer struct {
	Config *config.HTTPServer
	Server *http.Server
	Logger logger.Logger
}

func NewServer(opts *Opts) *HTTPServer {
	mux := http.NewServeMux()

	healthService := service.NewHealthService(&service.HealthServiceOpts{
		Database: opts.Database,
		Cache:    opts.Cache,
	})
	healthHandler := handler.NewHealthHandler(&handler.Opts{
		HealthService: healthService,
		Logger:        opts.Logger,
	})

	mux.HandleFunc("/livez", healthHandler.Livez)
	mux.HandleFunc("/readyz", healthHandler.Readyz)

	if opts.Metrics != nil {
		mux.Handle("/metrics", opts.Metrics.Handler())
	}

	return &HTTPServer{
		Config: opts.Config,
		Server: &http.Server{
			Addr:    opts.Config.URL,
			Handler: mux,
		},
		Logger: opts.Logger,
	}
}

func (h *HTTPServer) ServeListener(listener net.Listener) error {
	h.Logger.Info("HTTP server started", logger.Field{Key: "address", Value: listener.Addr().String()})
	if err := h.Server.Serve(listener); err != nil && err != http.ErrServerClosed {
		h.Logger.Error("HTTP server failed", logger.Field{Key: "error", Value: err.Error()})
		return err
	}
	return nil
}

func (h *HTTPServer) Serve() error {
	listener, err := net.Listen("tcp", h.Config.URL)
	if err != nil {
		h.Logger.Error("Failed to create HTTP listener",
			logger.Field{Key: "address", Value: h.Config.URL},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return err
	}

	return h.ServeListener(listener)
}
