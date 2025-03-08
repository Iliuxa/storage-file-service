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
	ErrUsersGettingError = errors.New("user not found")
)

type StorageRepository interface {
	FindAll(context.Context) ([]*StorageFile, error)
	InsertFile(context.Context, *StorageFile) (uint, error)
}
