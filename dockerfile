# Используем многоступенчатую сборку
FROM golang:1.24-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git

# Рабочая директория
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

# Финальный образ
FROM alpine:latest

# Установка инструментов
RUN apk add --no-cache postgresql-client

# Рабочая директория
WORKDIR /app


# Копируем бинарник
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY .env .
# Команда запуска
CMD ["./main"]