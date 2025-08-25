#!/bin/bash

# Path settings
DB_URL="postgres://postgres:postgres@localhost:5432/bugby
MIGRATIONS_DIR="internal/db/migrations"
SCHEMA_FILE="internal/db/schema.sql"
SQLC_DIR="internal/db/sqlc"

echo "⬆️  Running migrations with Goose..."
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up

echo "📄 Dumping schema to $SCHEMA_FILE..."
pg_dump "$DB_URL" --schema-only --no-owner --file="$SCHEMA_FILE"

echo "⚙️  Running sqlc generate..."
sqlc generate

echo "✅ Done: migrations + schema + sqlc"
