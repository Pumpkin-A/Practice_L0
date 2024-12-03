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
RUN go build -o /app/consumer ./cmd/main.go

# Финальный образ
FROM alpine
WORKDIR /app
# Копируем скомпилированное приложение
COPY --from=builder /app/consumer /app/consumer
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/.env /app/.env

# Экспорт переменных окружения для Kafka и PostgreSQL
# ENV KAFKA_BROKER=host.docker.internal:9092
# ENV POSTGRES_HOST=host.docker.internal
# ENV POSTGRES_PORT=5432
# ENV POSTGRES_USER=user_table_orders
# ENV POSTGRES_DB=ordersdb

EXPOSE 9090
# Запуск приложения
CMD ["./consumer"]
