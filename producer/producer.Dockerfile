# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /build
# Копируем и устанавливаем зависимости
ADD go.mod .
ADD go.sum .
RUN go mod download
# Копируем исходный код
COPY . .
RUN go build -o /app/producer ./cmd/main.go

# Финальный образ
FROM alpine
WORKDIR /app
# Копируем скомпилированное приложение
COPY --from=builder /app/producer /app/producer

# Запуск приложения
CMD ["./producer"]
