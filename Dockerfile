# Используем официальный образ Go для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию для сборки
WORKDIR /app

# Копируем go.mod и go.sum для оптимизации слоев
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем приложение
RUN GOOS=linux CGO_ENABLED=0 go build -o main ./cmd/main.go

# Используем легковесный образ для выполнения приложения
FROM alpine:latest

COPY --from=builder /app/main /app/main

# Устанавливаем рабочую директорию
WORKDIR /app

ENTRYPOINT ["./main"]
