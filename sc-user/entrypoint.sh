#!/bin/sh
set -eu

MIGRATIONS_DIR=${MIGRATIONS_DIR:-/app/internal/db/migration}
APP_BINARY=${APP_BINARY:-/usr/local/bin/general}
MIGRATE_BIN=${MIGRATE_BIN:-/usr/local/bin/migrate}

if [ -z "${DB_SOURCE:-}" ]; then
  echo "DB_SOURCE environment variable is required" >&2
  exit 1
fi

# Run database migrations
"${MIGRATE_BIN}" -path "${MIGRATIONS_DIR}" -database "${DB_SOURCE}" -verbose up

# Start the service
exec "${APP_BINARY}"
