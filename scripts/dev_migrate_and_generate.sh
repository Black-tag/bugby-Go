# #!/bin/bash

# # Path settings
# DB_URL="postgres://postgres:postgres@localhost:5432/bugby"
# MIGRATIONS_DIR="internal/db/migrations"
# SCHEMA_FILE="internal/db/schema.sql"
# SQLC_DIR="internal/db/sqlc"

# echo "⬆️  Running migrations with Goose..."
# goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up

# echo "📄 Dumping schema to $SCHEMA_FILE..."
# pg_dump "$DB_URL" --schema-only --no-owner --file="$SCHEMA_FILE"

# echo "⚙️  Running sqlc generate..."
# sqlc generate

# echo "✅ Done: migrations + schema + sqlc"


#!/usr/bin/env bash
set -euo pipefail

# Config
DB_HOST="${DB_HOST:-db}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${POSTGRES_USER:-postgres}"
DB_NAME="${POSTGRES_DB:-bugby}"
DB_PASS="${POSTGRES_PASSWORD:-postgres}"
# Fallback DB_URL if not provided
: "${DB_URL:=postgres://postgres:postgres@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-internal/db/migrations}"
SCHEMA_FILE="${SCHEMA_FILE:-internal/db/schema.sql}"
SQLC_DIR="${SQLC_DIR:-internal/db/sqlc}"
MAX_WAIT="${MAX_WAIT_SECONDS:-60}"

echo "🔌 Using database URL: ${DB_URL}"
echo "🔄 Waiting for Postgres to be ready (host=${DB_HOST} port=${DB_PORT})..."

i=0
while true; do
  if command -v pg_isready >/dev/null 2>&1; then
    if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1; then
      echo "✅ pg_isready: Postgres is ready"
      break
    fi
  else
    # fallback TCP check
    if command -v nc >/dev/null 2>&1 && nc -z "$DB_HOST" "$DB_PORT" >/dev/null 2>&1; then
      echo "✅ tcp port open for Postgres"
      break
    fi
  fi

  i=$((i+1))
  if [ "$i" -ge "$MAX_WAIT" ]; then
    echo "❌ Timeout waiting for Postgres after ${MAX_WAIT}s"
    exit 1
  fi
  printf "Waiting... %ds\r" "$i"
  sleep 1
done

# Run migrations (goose) if available
if command -v goose >/dev/null 2>&1; then
  echo "⬆️  Running migrations with Goose..."
  # goose accepts a postgres URL without shell expansion issues
  goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up
else
  echo "⚠️  goose not found in image, skipping migrations (install goose in Dockerfile if you need migrations here)"
fi

# Dump schema if pg_dump exists
if command -v pg_dump >/dev/null 2>&1; then
  echo "📄 Dumping schema to $SCHEMA_FILE..."
  # pg_dump may require env PGPASSWORD for non-interactive auth
  export PGPASSWORD="${DB_PASS}"
  pg_dump "$DB_URL" --schema-only --no-owner --file="$SCHEMA_FILE" || echo "⚠️  pg_dump failed (maybe file permissions)"
  unset PGPASSWORD
else
  echo "⚠️  pg_dump not found, skipping schema dump"
fi

# Run sqlc generate if sqlc exists
if command -v sqlc >/dev/null 2>&1; then
  echo "⚙️  Running sqlc generate..."
  sqlc generate
else
  echo "⚠️  sqlc not found, skipping sqlc generation"
fi

echo "✅ bootstrap complete — starting application"
/app/bugby-server
