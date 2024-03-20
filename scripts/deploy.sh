#!/usr/bin/env bash

export REPMGR_PARTNER_NODES=pg-1,pg-2,ubuntu-server-1-pg-1,ubuntu-server-1-pg-2
export PGPOOL_BACKEND_NODES=1:pg-1:5432,2:pg-2:5432,3:192.168.188.132:5432,4:192.168.188.132:5433
export PGPOOL_BACKEND_APPLICATION_NAMES=pg-1,pg-2,ubuntu-server-1-pg-1,ubuntu-server-1-pg-2

docker compose -f docker-compose.postgresql.yml up -d