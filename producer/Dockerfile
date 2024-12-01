# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию для сборки
WORKDIR /app

# Копируем файлы для управления зависимостями
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект
COPY . .


# Меняем рабочую директорию на /app/cmd, где находится main.go
WORKDIR /app/producer/cmd

# Сборка приложения
RUN go build -o /producer .

# Финальный минимальный образ
FROM debian:bullseye-slim

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированное приложение
COPY --from=builder /producer .

# Устанавливаем переменные окружения (если требуется)
ENV KAFKA_BROKER=kafka:9092

# Запуск приложения
CMD ["./producer"]
