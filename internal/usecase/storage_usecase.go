package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"storage-file-service/internal/domain"
	proto "storage-file-service/proto/gen"
	"time"
)

type StorageUsecase interface {
	GetFileList(context.Context) ([]*proto.FileInfo, error)
	UploadFile(context.Context, *bytes.Buffer, string, string) (uint, error)
	DownloadFile(context.Context, uint) (*proto.FileResponse_Info, string, error)
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
		if errors.Is(err, domain.ErrFileNotFound) {
			s.log.Error("file mot found ", slog.StringValue(err.Error()))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}

		log.Error("Failed to getting users", slog.StringValue(err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fileInfos := make([]*proto.FileInfo, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, &proto.FileInfo{
			FileId:    uint64(file.Id),
			Filename:  file.Name,
			CreatedAt: file.InsertDate,
			UpdatedAt: file.UpdateDate.String,
		})
	}

	return fileInfos, nil
}

func (s *storageUsecase) UploadFile(ctx context.Context, fileBuffer *bytes.Buffer, fileName, fileHash string) (uint, error) {
	const op = "StorageUsecase.SaveFile"
	log := s.log.With(
		slog.String("op", op),
	)
	log.Info("Saving file")

	storageFilePath, err := s.saveFile(*fileBuffer, fileName, fileHash)
	if err != nil {
		os.Remove(storageFilePath)
		s.log.Error("Failed to create file", slog.StringValue(err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	storageFile := domain.StorageFile{
		Name:       fileName,
		InsertDate: time.Now().Format("2006-01-02 15:04:05"),
		UpdateDate: sql.NullString{
			String: time.Now().Format("2006-01-02 15:04:05"),
			Valid:  true,
		},
		FilePath: storageFilePath,
		FileHash: fileHash,
	}

	fileId, err := s.storageRepository.InsertFile(ctx, &storageFile)
	if err != nil {
		os.Remove(storageFilePath)

		s.log.Error("Failed to create storage file", slog.StringValue(err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return fileId, nil
}

func (s *storageUsecase) DownloadFile(ctx context.Context, fileId uint) (*proto.FileResponse_Info, string, error) {
	const op = "StorageUsecase.DownloadFile"
	log := s.log.With(
		slog.String("op", op),
	)
	log.Info("Downloading file")

	storageFile, err := s.storageRepository.Find(ctx, fileId)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
	}

	file, err := os.Open(storageFile.FilePath)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
	}
	defer file.Close()

	if err = s.checkingFileHash(storageFile.FileHash, file); err != nil {
		return nil, "", fmt.Errorf("%s: %w", "checing file", domain.ErrFileIsDamaged)
	}

	fileInfo := &proto.FileResponse_Info{
		Info: &proto.ImageInfo{
			FileName: storageFile.Name,
			FileHash: storageFile.FileHash,
		},
	}

	return fileInfo, storageFile.FilePath, nil
}

func (s *storageUsecase) saveFile(fileBuffer bytes.Buffer, fileName, fileHash string) (string, error) {
	const storageDir = "/storage/files/"

	date := time.Now().Format("2006/01/02/")
	exePath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	storageFilePath := exePath + storageDir + date

	if err := os.MkdirAll(storageFilePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	storageFilePath += s.getRandomString(10) + "_" + fileName
	file, err := os.Create(storageFilePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}
	defer file.Close()

	_, err = fileBuffer.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("cannot seek file: %w", err)
	}

	if err = s.checkingFileHash(fileHash, file); err != nil {
		return "", fmt.Errorf("%s: %w", "checing file", domain.ErrFileIsDamaged)
	}

	return storageFilePath, nil
}

func (s *storageUsecase) getRandomString(size uint) string {
	const charset = "0123456789"

	result := make([]byte, size)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func (s *storageUsecase) checkingFileHash(fileHash string, file io.Reader) error {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		s.log.Error("Failed to get hash")
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	if fileHash != hash {
		s.log.Error("Failed to check file hash")
	}
	return nil
}
