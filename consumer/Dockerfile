# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Меняем рабочую директорию на /app/cmd, где находится main.go
WORKDIR /app/consumer/cmd

# Сборка приложения
RUN go build -o consumer .

# Финальный образ
FROM debian:bullseye-slim

WORKDIR /root/

# Копируем скомпилированное приложение
COPY --from=builder /app/consumer .

# Экспорт переменных окружения для Kafka и PostgreSQL
ENV KAFKA_BROKER=kafka:9092
ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=practice_db

# Запуск приложения
CMD ["./consumer"]
