#!/bin/sh
set -eu

MIGRATIONS_DIR=${MIGRATIONS_DIR:-/app/internal/db/migration}
APP_BINARY=${APP_BINARY:-/usr/local/bin/general}
MIGRATE_BIN=${MIGRATE_BIN:-/usr/local/bin/migrate}

DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-secret}
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-sc_db}
DB_SSL_MODE=${DB_SSL_MODE:-disable}

DB_SOURCE=${DB_SOURCE:-"postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"}

if [ -z "${DB_SOURCE:-}" ]; then
  echo "DB_SOURCE environment variable is required" >&2
  exit 1
fi

# Run database migrations
"${MIGRATE_BIN}" -path "${MIGRATIONS_DIR}" -database "${DB_SOURCE}" -verbose up

# Start the service
exec "${APP_BINARY}"
