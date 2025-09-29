#!/busybox/sh
set -euo pipefail

MIGRATION_DIR="${MIGRATION_DIR:-/app/internal/db/migration}"
DB_SOURCE="${DB_SOURCE:?DB_SOURCE environment variable is required}"

/bin/migrate -path "${MIGRATION_DIR}" -database "${DB_SOURCE}" up

exec /app/auth-service
