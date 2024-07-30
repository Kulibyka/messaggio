#!/bin/bash

set -e

host="$1"
shift
cmd="$@"

until psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$host:5432/postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - creating database if not exists"

psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$host:5432/postgres" <<-EOSQL
 CREATE DATABASE $POSTGRES_DB;
EOSQL
#psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$host:5432/postgres" <<-EOSQL
#  SELECT 'CREATE DATABASE $POSTGRES_DB'
#  WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$POSTGRES_DB')\gexec;
#EOSQL

>&2 echo "Database created or already existed - running migrations"

/app/migrator

>&2 echo "Migrations applied - executing command"

exec $cmd
