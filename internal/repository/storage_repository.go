package repository

import (
	"context"
	"database/sql"
	"fmt"
	"storage-file-service/internal/domain"
)

type storageRepository struct {
	db *sql.DB
}

func NewStorageRepository(storagePath string) (domain.StorageRepository, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storageRepository{db: db}, nil
}

func (r *storageRepository) FindAll(ctx context.Context) ([]*domain.StorageFile, error) {
	const op = "repository.sqlite.findAll"

	stmt, err := r.db.Prepare("SELECT * FROM storage_file")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var storageFiles []*domain.StorageFile

	for rows.Next() {
		var storageFile domain.StorageFile
		if err := rows.Scan(&storageFile.Id, &storageFile.Name, &storageFile.InsertDate, &storageFile.UpdateDate, &storageFile.FilePath, &storageFile.FileHash); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		storageFiles = append(storageFiles, &storageFile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storageFiles, nil
}
