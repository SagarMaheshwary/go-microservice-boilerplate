package server

import (
	"context"

	example "github.com/sagarmaheshwary/go-microservice-boilerplate/proto"
)

type ExampleServer struct {
	example.ExampleServer
}

func (e *ExampleServer) Hello(ctx context.Context, in *example.HelloRequest) (*example.HelloResponse, error) {
	return &example.HelloResponse{
		Message: "Hello World",
	}, nil
}
