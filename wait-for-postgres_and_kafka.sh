#!/bin/bash

set -e

host="$1"
shift
cmd="$@"

# Ожидание подъема Postgres
until psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$host:5432/postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - creating database if not exists"

psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$host:5432/postgres" <<-EOSQL
 CREATE DATABASE $POSTGRES_DB;
EOSQL

>&2 echo "Database created or already existed - running migrations"

/app/migrator

>&2 echo "Migrations applied - executing command"

# Ожидание подъема Kafka
KAFKA_HOST="kafka"   # Это значение нужно указать в соответствии с тем, как называется ваш Kafka контейнер в Docker Compose
KAFKA_PORT=9092

until echo > /dev/tcp/$KAFKA_HOST/$KAFKA_PORT; do
  >&2 echo "Kafka is unavailable - sleeping"
  sleep 1
done

>&2 echo "Kafka is up - executing command"

exec $cmd
