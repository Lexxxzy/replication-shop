#!/usr/bin/env bash

SERVER_IP=$1

export REPMGR_PARTNER_NODES=pg-1,pg-2
export PGPOOL_BACKEND_NODES=1:pg-1:5432,2:pg-2:5432,3:${SERVER_IP}:5432,4:${SERVER_IP}:5433

docker compose -f docker-compose.postgresql.yml up -d