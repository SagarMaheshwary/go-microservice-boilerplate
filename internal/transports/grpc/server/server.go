package server

import (
	"net"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server/handler"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/transports/grpc/server/interceptor"
	helloworld "github.com/sagarmaheshwary/go-microservice-boilerplate/proto/hello_world"
	"google.golang.org/grpc"
)

type Opts struct {
	Config   *config.GRPCServer
	Logger   logger.Logger
	Database database.DatabaseService
	Cache    cache.CacheService
}

type GRPCServer struct {
	Server *grpc.Server
	Config *config.GRPCServer
	Logger logger.Logger
}

func NewServer(opts *Opts) *GRPCServer {
	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptor.LoggerInterceptor(opts.Logger)))
	helloworld.RegisterGreeterServer(srv, handler.NewGreeterServer(
		service.NewUserService(&service.UserServiceOpts{
			Database: opts.Database,
			Cache:    opts.Cache,
		}),
	))

	return &GRPCServer{
		Server: srv,
		Config: opts.Config,
		Logger: opts.Logger,
	}
}

func (s *GRPCServer) ServeListener(listener net.Listener) error {
	s.Logger.Info("gRPC server started", logger.Field{Key: "address", Value: listener.Addr().String()})
	if err := s.Server.Serve(listener); err != nil {
		s.Logger.Error("gRPC server failed", logger.Field{Key: "error", Value: err.Error()})
		return err
	}
	return nil
}

func (s *GRPCServer) Serve() error {
	url := s.Config.URL
	listener, err := net.Listen("tcp", url)
	if err != nil {
		s.Logger.Error("Failed to create tcp listener",
			logger.Field{Key: "address", Value: url},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return err
	}
	return s.ServeListener(listener)
}
