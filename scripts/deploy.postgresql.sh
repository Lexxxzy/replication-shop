#!/usr/bin/env bash

set -e

SERVER_IP=$1

if [ -z "${SERVER_IP}" ]; then
  echo "variable SERVER_IP is empty"
  exit 1
fi

PGPOOL_BACKEND_NODES=1:pg-1:5432,2:pg-2:5432,3:${SERVER_IP}:5433,4:${SERVER_IP}:5435 \
docker compose -f docker-compose.postgresql.yml up -d