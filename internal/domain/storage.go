package domain

import (
	"context"
	"errors"
)

type StorageFile struct {
	Id         int
	Name       string
	InsertDate string
	UpdateDate string
	FilePath   string
	FileHash   string
}

var (
	ErrUsersGettingError = errors.New("user not found")
)

type StorageRepository interface {
	FindAll(context.Context) ([]*StorageFile, error)
}
