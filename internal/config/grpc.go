package config

import (
	"context"
	"fmt"
	"grpc-server-1/internal/usecase"
	"grpc-server-1/protogen/user"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GrpcAdapter struct {
	UserUseCase *usecase.UserUseCase
	user.UserServiceServer
	server   *grpc.Server
	grpcPort int
}

func NewGrpcAdapter(userUseCase *usecase.UserUseCase, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		UserUseCase: userUseCase,
		grpcPort:    grpcPort,
	}
}

func (a *GrpcAdapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d : %v\n", a.grpcPort, err)
	}

	log.Printf("Server listening on port %d\n", a.grpcPort)

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	user.RegisterUserServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve gRPC on port %d : %v\n", a.grpcPort, err)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}

func (a *GrpcAdapter) GetUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	return a.UserUseCase.GetUser(ctx, req)
}

func (a *GrpcAdapter) CreateUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	return a.UserUseCase.CreateUser(ctx, req)
}

func (a *GrpcAdapter) SearchUser(ctx context.Context, req *user.UserSearchRequest) (*user.PaginatedUserResponse, error) {
	return a.UserUseCase.SearchUser(ctx, req)
}

func (a *GrpcAdapter) UpdateUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	return a.UserUseCase.UpdateUser(ctx, req)
}
