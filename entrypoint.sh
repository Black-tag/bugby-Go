#!/bin/sh
set -e


# Use explicit Railway variables if DATABASE_URL not set
DB_URL=${DATABASE_URL:-"postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=require"}

DB_URL=${DATABASE_URL:-"postgresql://postgres:postgres@db:5432/bugby?sslmode=disable"}
echo "🔌 Using database URL: ${DB_URL}"
echo "🔄 Waiting for Postgres to be ready..."
until psql "postgresql://postgres:postgres@db:5432/bugby?sslmode=disable" -c '\q' > /dev/null 2>&1; do
  sleep 1
done

echo "✅ Running migrations with Goose..."
goose -dir ./internal/db/migrations postgres "postgres://postgres:postgres@db:5432/bugby?sslmode=disable" up

echo "✅ Running migrations with Goose..."
goose -dir ./internal/db/migrations postgres "${DB_URL}" up

echo "🚀 Starting Go app..."
./bugby
