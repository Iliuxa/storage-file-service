# Используем официальный образ Go для сборки
FROM golang:1.23.4-alpine3.21

WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Компилируем API и миграции
RUN go build -o app ./cmd/app/main.go
RUN go build -o migrator ./cmd/migrator/main.go

# Делаем файлы исполняемыми
RUN chmod +x /app/migrator /app/app

# Запускаем миграции, а затем основное приложение
CMD ["sh", "-c", "./migrator --storage-path=./storage/storage_file.db --migrations-path=./deployment/migrations && ./app"]
