package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"storage-file-service/internal/domain"
	proto "storage-file-service/proto/gen"
)

type StorageUsecase interface {
	GetFileList(context.Context) ([]*proto.FileInfo, error)
}

type storageUsecase struct {
	storageRepository domain.StorageRepository
	log               *slog.Logger
}

func NewStorageUsecase(storageRepository domain.StorageRepository, log *slog.Logger) StorageUsecase {
	return &storageUsecase{
		storageRepository: storageRepository,
		log:               log,
	}
}

func (s *storageUsecase) GetFileList(ctx context.Context) ([]*proto.FileInfo, error) {
	const op = "StorageUsecase.GetFileList"
	log := s.log.With(
		slog.String("op", op),
	)
	log.Info("getting files info")

	files, err := s.storageRepository.FindAll(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrUsersGettingError) {
			s.log.Error("Failed to getting users", slog.StringValue(err.Error()))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrUsersGettingError)
		}

		log.Error("Failed to getting users", slog.StringValue(err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fileInfos := make([]*proto.FileInfo, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, &proto.FileInfo{
			Filename:  file.Name,
			CreatedAt: file.InsertDate,
			UpdatedAt: file.UpdateDate,
		})
	}

	return fileInfos, nil
}
