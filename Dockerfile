# Используем официальный образ Go
FROM golang:1.25.5-alpine AS builder

# Устанавливаем зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем зависимости для работы приложения
RUN apk --no-cache add ca-certificates

# Создаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл из builder образа
COPY --from=builder /app/main .

# Копируем статические файлы
COPY --from=builder /app/static ./static

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]