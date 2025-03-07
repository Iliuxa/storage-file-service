package app

import (
	"log/slog"
	grpcapp "storage-file-service/internal/app/grpc"
	"storage-file-service/internal/repository"
	"storage-file-service/internal/usecase"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	storageRepository, err := repository.NewStorageRepository(storagePath)
	if err != nil {
		panic(err)
	}
	storageUsecase := usecase.NewStorageUsecase(storageRepository, log)
	grpcApp := grpcapp.New(log, storageUsecase, grpcPort)

	return &App{GRPCServer: grpcApp}
}
