#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

echo "Waiting for Postgres with the following details:"
echo "Host: $host"
echo "User: $DATABASE_USER"
echo "Password: $DATABASE_PASSWORD"
echo "DSN: $DATABASE_DSN"

until PGPASSWORD=$DATABASE_PASSWORD psql -h "$host" -U "$DATABASE_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
