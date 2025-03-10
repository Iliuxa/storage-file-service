# storage-file-service

Выполнение миграции
```sh
go run ./cmd/migrator --storage-path=./storage/storage_file.db --migrations-path=./deployment/migrations
```

Выполнение миграции для тестов
```sh
go run ./cmd/migrator --storage-path=./storage/storage_file.db --migrations-path=./tests/migrations --migrations-table=migrations_test
```
--------------------------

Разворачивание в контейнере
```sh
docker-compose up --build
```
Для запуска тестов зайти в контейнер go-storage-app или в рабочей директории при запуске без контейнера
```sh
go test ./tests/
```