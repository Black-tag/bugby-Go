#!/bin/sh
set -e





DB_URL="${DB_URL}"
echo "🔌 Using database URL: ${DB_URL}"
echo "🔄 Waiting for Postgres to be ready..."
until psql "$DB_URL" -c '\q' > /dev/null 2>&1; do
  sleep 1
done

echo "✅ Running migrations with Goose..."
goose -dir ./internal/db/migrations postgres "postgres://postgres:postgres@db:5432/bugby?sslmode=disable" up

echo "✅ Running migrations with Goose..."
goose -dir ./internal/db/migrations postgres "$DB_URL" up

echo "🚀 Starting Go app..."
./bugby
