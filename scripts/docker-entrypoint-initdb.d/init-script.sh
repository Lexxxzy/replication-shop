#!/bin/sh
set -e

DB_NAME="${POSTGRESQL_DATABASE:-postgres}"

result=$(PGPASSWORD="$POSTGRESQL_PASSWORD" psql -U "$POSTGRESQL_USERNAME" -h localhost -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'")
if [ -z "$result" ]; then
    echo "[INIT] Database $DB_NAME does not exist. Importing dump..."
    PGPASSWORD="$POSTGRESQL_PASSWORD" psql -U "$POSTGRESQL_USERNAME" -d "$DB_NAME" -f /docker-entrypoint-initdb.d/dump.sql
else
    echo "[INIT] Database $DB_NAME already exists. Skipping import."
fi
