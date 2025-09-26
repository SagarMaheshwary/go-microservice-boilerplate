package handler

import (
	"context"

	helloworld "github.com/sagarmaheshwary/go-microservice-boilerplate/proto/hello_world"
)

type GreeterServer struct {
	helloworld.GreeterServer
}

func (e *GreeterServer) SayHello(ctx context.Context, in *helloworld.SayHelloRequest) (*helloworld.SayHelloResponse, error) {
	return &helloworld.SayHelloResponse{
		Message: "Hello, World!",
	}, nil
}
