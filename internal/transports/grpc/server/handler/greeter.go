package handler

import (
	"context"
	"fmt"

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
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &helloworld.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", user.Name),
		User: &helloworld.User{
			Id:    int64(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
