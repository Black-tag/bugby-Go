#!/bin/sh
set -e

echo "🔄 Waiting for Postgres to be ready..."
until psql "postgresql://postgres:postgres@db:5432/bugby?sslmode=disable" -c '\q' > /dev/null 2>&1; do
  sleep 1
done

echo "✅ Running migrations with Goose..."
goose -dir ./internal/db/migrations postgres "postgres://postgres:postgres@db:5432/bugby?sslmode=disable" up

echo "🚀 Starting Go app..."
./bugby
