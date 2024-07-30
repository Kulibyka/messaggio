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
FROM golang:1.22

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



## Используем официальный образ Golang для сборки приложения
#FROM golang:1.22
#
## Устанавливаем рабочую директорию внутри контейнера
#WORKDIR /app
#
## Копируем go.mod и go.sum для установки зависимостей
#COPY go.mod go.sum ./
#
#RUN apt-get update
#RUN apt-get -y install postgresql-client
#
## Устанавливаем зависимости
#RUN go mod download
#
## Копируем исходный код
#COPY ./ ./
#
## Собираем приложение
#RUN CGO_ENABLED=0 GOOS=linux go build -o ./messaggio ./cmd/messageSrv
#RUN CGO_ENABLED=0 GOOS=linux go build -o ./migrate ./cmd/migrator
#
## Копируем конфигурационные файлы и скрипт
#COPY ./config/local.yaml ./config.yaml
#COPY ./wait-for-postgres.sh ./wait-for-postgres.sh
#COPY ./migrations ./migrations
#
#RUN chmod +x ./wait-for-postgres.sh
#
## Проверка наличия файлов
#RUN ls -la
#
#EXPOSE 8080
#
## ENTRYPOINT для выполнения миграций и запуска приложения
#ENTRYPOINT ["./migrate", "./wait-for-postgres.sh", "db", "./messaggio"]

## Используем официальный образ Golang для сборки приложения
#FROM golang:1.22.0 AS builder
#
## Устанавливаем рабочую директорию внутри контейнера
#WORKDIR /app
#
## Копируем go.mod и go.sum для установки зависимостей
#COPY go.mod go.sum ./
#
## Устанавливаем зависимости
#RUN go mod download
#
## Копируем исходный код
#COPY . .
#
## Собираем приложение
#RUN CGO_ENABLED=0 go build -o /app/messaggio ./cmd/messageSrv
##RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate
### Создаем минимальный образ для запуска нашего приложения
##FROM alpine
#
## Устанавливаем рабочую директорию внутри контейнера
#WORKDIR /root/
#
## Копируем скомпилированное приложение из предыдущего этапа
##COPY --from=builder /app/messaggio .
##COPY config/local.yaml local.yaml
#COPY ./wait-for-postgres.sh .
#
## Копируем конфигурационные файлы и скрипты
#
#
#RUN apt-get update
#RUN apt-get -y install postgresql-client
##
## Устанавливаем права на выполнение для скрипта
#RUN chmod +x ./wait-for-postgres.sh
#
#EXPOSE 8080
#
#ENTRYPOINT ["./wait-for-postgres.sh", "db", "./app/messaggio"]