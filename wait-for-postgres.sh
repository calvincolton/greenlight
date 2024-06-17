#!/bin/sh
set -e

# Use DATABASE_NAME if set, otherwise default to POSTGRES_DB
DB_NAME=${DATABASE_NAME:-$POSTGRES_DB}

host="$1"
shift
cmd="$@"

echo "Waiting for Postgres with the following details:"
echo "Host: $host"
echo "DATABASE_USER: $DATABASE_USER"
echo "DATABASE_PASSWORD: $DATABASE_PASSWORD"
echo "DB_NAME: $DB_NAME"
echo "DATABASE_DSN: $DATABASE_DSN"

until PGPASSWORD=$DATABASE_PASSWORD psql -h "$host" -U "$DATABASE_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
