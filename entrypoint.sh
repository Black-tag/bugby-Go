#!/bin/sh
set -e





# DB_URL="${DB_URL}"
echo "🔌 Using database URL: ${DB_URL}"
echo "🔄 Waiting for Postgres to be ready..."
until pg_isready "$DB_URL" -c '\q' > /dev/null 2>&1; do
  echo "postgres is unavailible sleepin...."
  sleep 1
done
# until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
#   sleep 2
# done

echo "✅ Running migrations with Goose..."
goose -dir ./internal/db/migrations postgres "$DB_URL" up

# echo "✅ Running migrations with Goose..."
# goose -dir ./internal/db/migrations postgres "$DB_URL" up

echo "🚀 Starting Go app..."
./bugby
