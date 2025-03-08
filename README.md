start migration
```sh
go run ./cmd/migrator --storage-path=./storage/storage_file.db --migrations-path=./deployment/migrations
```

start test migration
```sh
go run ./cmd/migrator --storage-path=./storage/storage_file.db --migrations-path=./tests/migrations --migrations-table=migrations_test
```