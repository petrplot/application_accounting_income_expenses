# Этап 1: Сборка
FROM golang:latest AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
#COPY go.mod go.sum ./
COPY . .
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник (статическая линковка для работы в alpine/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/app/main.go

# Этап 2: Финальный образ
FROM alpine:latest

# Добавляем сертификаты для работы с HTTPS (нужно для API и внешних запросов)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем только скомпилированный файл из первого этапа
COPY --from=builder /main .
# Копируем .env (если он нужен при запуске внутри контейнера)
COPY .env .

RUN chmod +x ./main
# Запуск приложения
CMD ["./main"]
