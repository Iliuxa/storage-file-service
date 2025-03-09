package domain

import (
	"context"
	"database/sql"
	"errors"
)

type StorageFile struct {
	Id         int
	Name       string
	InsertDate string
	UpdateDate sql.NullString
	DeleteDate sql.NullString
	FilePath   string
	FileHash   string
}

var (
	ErrFileGettingError   = errors.New("file not found")
	ErrFileIsDamagedError = errors.New("the file is damaged")
)

type StorageRepository interface {
	FindAll(context.Context) ([]*StorageFile, error)
	InsertFile(context.Context, *StorageFile) (uint, error)
	Find(ctx context.Context, id uint) (*StorageFile, error)
}
