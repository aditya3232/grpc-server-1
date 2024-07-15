package main

import (
	"grpc-server-1/internal/config"
	"grpc-server-1/internal/repository"
	"grpc-server-1/internal/usecase"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	logger.Info("Starting grpc service")

	userRepository := repository.NewUserRepository(logger)
	userUsecase := usecase.NewUserUseCase(db, logger, userRepository)
	grpcAdapter := config.NewGrpcAdapter(userUsecase, 9090)
	logger.Info("grpc service running")

	grpcAdapter.Run()
}
