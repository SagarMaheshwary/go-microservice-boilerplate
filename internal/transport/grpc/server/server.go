package server

import (
	"net"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
	example "github.com/sagarmaheshwary/go-microservice-boilerplate/proto"
	"google.golang.org/grpc"
)

type Opts struct {
	Config   *config.GRPCServer
	Logger   logger.Logger
	Database database.DatabaseService
}

type GRPCServer struct {
	Server   *grpc.Server
	Config   *config.GRPCServer
	Logger   logger.Logger
	Database database.DatabaseService
}

func NewServer(opts *Opts) *GRPCServer {
	srv := grpc.NewServer()
	example.RegisterExampleServer(srv, &ExampleServer{})

	return &GRPCServer{
		Server:   srv,
		Config:   opts.Config,
		Logger:   opts.Logger,
		Database: opts.Database,
	}
}

func (s *GRPCServer) ServeListener(listener net.Listener) error {
	s.Logger.Info("gRPC server started on %q", listener.Addr().String())
	if err := s.Server.Serve(listener); err != nil {
		s.Logger.Error("gRPC server failed: %v", err)
		return err
	}
	return nil
}

func (s *GRPCServer) Serve() error {
	url := s.Config.URL
	listener, err := net.Listen("tcp", url)
	if err != nil {
		s.Logger.Error("Failed to create tcp listener on %q: %v", url, err)
		return err
	}
	return s.ServeListener(listener)
}
