# Используем официальный образ Golang для сборки приложения
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код
COPY ./ ./

# Собираем приложение
RUN  go build -o /app/messaggio ./cmd/messageSrv
RUN  go build -o /app/migrator ./cmd/migrator

# Создаем финальный образ для запуска нашего приложения
FROM ubuntu:22.04

RUN apt-get update
RUN apt-get -y install postgresql-client
# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/messaggio .
COPY --from=builder /app/migrator .

# Копируем конфигурационные файлы
COPY config/local.yaml .
COPY wait-for-postgres.sh .

# Копируем миграции
COPY migrations ./migrations

RUN chmod +x ./wait-for-postgres.sh

EXPOSE 8080
ENTRYPOINT ["./wait-for-postgres.sh", "db", "./messaggio"]
