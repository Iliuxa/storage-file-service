package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
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

	stmt, err := r.db.Prepare("select id, file_name, insert_date, update_date, delete_date, file_path, file_hash from storage_file where delete_date is null")
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
		if err := rows.Scan(&storageFile.Id, &storageFile.Name, &storageFile.InsertDate, &storageFile.UpdateDate, &storageFile.DeleteDate, &storageFile.FilePath, &storageFile.FileHash); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		storageFiles = append(storageFiles, &storageFile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storageFiles, nil
}

func (r *storageRepository) InsertFile(ctx context.Context, file *domain.StorageFile) (uint, error) {
	const op = "repository.sqlite.insertFile"

	stmt, err := r.db.Prepare(`
		INSERT INTO storage_file (file_name, insert_date, update_date, file_path, file_hash)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, file.Name, file.InsertDate, file.UpdateDate, file.FilePath, file.FileHash)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
