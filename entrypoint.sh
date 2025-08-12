#!/bin/sh
set -e

# Debug: Print all DB-related variables (Alpine compatible)
echo "🔍 Environment Variables:"
env | awk -F= '/DATABASE|PG/ {print}'

# Priority: DATABASE_URL > Railway PG* vars > Local fallback
if [ -n "$DATABASE_URL" ]; then
    DB_URL="$DATABASE_URL"
elif [ -n "$PGHOST" ]; then
    DB_URL="postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=require"
else
    DB_URL="postgresql://postgres:postgres@localhost:5432/bugby?sslmode=disable"
    echo "⚠️  Using LOCAL database - for DEVELOPMENT only"
fi


# DB_URL=${DATABASE_URL:-"postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=require"}

# DB_URL=${DATABASE_URL:-"postgresql://postgres:postgres@db:5432/bugby?sslmode=disable"}
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
