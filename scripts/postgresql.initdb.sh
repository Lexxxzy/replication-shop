#!/usr/bin/env bash

set -e

export PGPASSWORD="$POSTGRESQL_PASSWORD"

table_count=$(psql -U "$POSTGRESQL_USERNAME" -d "$POSTGRESQL_DATABASE" \
-tAc "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'")

if [ "$table_count" -lt 1 ]; then
    echo "Database $POSTGRESQL_DATABASE does not exist. Importing dump..."
    psql -U "$POSTGRESQL_USERNAME" -d "$POSTGRESQL_DATABASE" -f /statements/ddl/postgresql.sql
    psql -U "$POSTGRESQL_USERNAME" -d "$POSTGRESQL_DATABASE" -f /statements/dml/postgresql.sql
else
    echo "Database $POSTGRESQL_DATABASE already exists. Skipping import."
fi
