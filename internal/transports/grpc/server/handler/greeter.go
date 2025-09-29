package handler

import (
	"context"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/service"
	helloworld "github.com/sagarmaheshwary/go-microservice-boilerplate/proto/hello_world"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreeterServer struct {
	helloworld.GreeterServer
	userService service.UserService
}

func NewGreeterServer(userService service.UserService) *GreeterServer {
	return &GreeterServer{userService: userService}
}

func (g *GreeterServer) SayHello(ctx context.Context, in *helloworld.SayHelloRequest) (*helloworld.SayHelloResponse, error) {
	user, err := g.userService.FindByID(ctx, uint(in.UserId))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &helloworld.SayHelloResponse{
		Message: "Hello, World!",
		User: &helloworld.User{
			Id:    int64(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
