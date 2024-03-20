#!/usr/bin/env bash

set -e

SERVER_IP=$1

if [ -z "${SERVER_IP}" ]; then
  echo "variable SERVER_IP is empty"
  exit 1
fi

docker compose -f docker-compose.postgresql.yml -f docker-compose.redis.yml -f docker-compose.app.yml down

PGPOOL_BACKEND_NODES=1:pg-1:5432,2:pg-2:5432,3:${SERVER_IP}:5432,4:${SERVER_IP}:5433 \
docker compose -f docker-compose.postgresql.yml up -d

docker compose -f docker-compose.redis.yml -f docker-compose.app.yml up -d