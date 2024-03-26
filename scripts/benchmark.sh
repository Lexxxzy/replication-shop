#!/usr/bin/env bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

source "$SCRIPT_DIR"/../.env

BRANCHES=${BRANCHES}

SLEEP_SECONDS=10

down_service() {
  local files
  files=$(find "$SCRIPT_DIR"/../ -type f -name 'docker-compose*.yml')

  for file in $files; do
    printf "${RED}Stopping %s containers${NC}\n" "$file"
    docker compose -f "$file" down
  done
  if [ "$(docker volume ls -q | wc -l)" -ge 1 ]; then
    printf "${RED}Remove all volumes${NC}\n"
    docker volume rm $(docker volume ls -q)
  fi
}

down_service

for branch_name in "${BRANCHES[@]}"; do
  printf "${GREEN}Branch: %s${NC}\n" "$branch_name"

  printf "${GREEN}Building images${NC}\n"
  docker build "https://github.com/Lexxxzy/replication-shop.git#$branch_name" \
    -t "replication-shop-app-$branch_name"

  printf "${GREEN}Start services${NC}\n"

  compose_files+=(-f docker-compose.postgresql.yml -f docker-compose.cassandra.yml -f docker-compose.redis.yml)
  case "$branch_name" in
  postgresql)
    export POSTGRESQL_ENABLED=true
    export CASSANDRA_ENABLED=false
    export REDIS_ENABLED=false
    ;;
  postgresql-cassandra)
    export POSTGRESQL_ENABLED=true
    export CASSANDRA_ENABLED=true
    export REDIS_ENABLED=false
    ;;
  postgresql-cassandra-redis)
    export POSTGRESQL_ENABLED=true
    export CASSANDRA_ENABLED=true
    export REDIS_ENABLED=true
    ;;
  esac

  printf "${GREEN}Wait for services to start${NC}\n"

  docker compose "${compose_files[@]}" up -d

  BRANCH_WITH_PREFIX=-${branch_name} \
    POSTGRESQL_ENABLED="$POSTGRESQL_ENABLED" \
    CASSANDRA_ENABLED="$CASSANDRA_ENABLED" \
    REDIS_ENABLED="$REDIS_ENABLED" \
    docker compose -f docker-compose.app.yml up -d

  printf "${GREEN}Wait for docker-compose.app.yml to start${NC}\n"
  sleep $((SLEEP_SECONDS + 80))

  printf "${GREEN}Start benchmark${NC}\n"
  BRANCH=${branch_name} \
    BRANCH_WITH_PREFIX=_${branch_name} \
    docker compose -f "$SCRIPT_DIR"/../docker-compose.benchmark.yml up

  down_service
done
