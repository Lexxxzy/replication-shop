#!/usr/bin/env bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

source "$SCRIPT_DIR"/../.env

BENCHMARK_INSTANCES_NUM=${BENCHMARK_INSTANCES_NUM:-6}
BENCHMARK_OUTPUT_DIR=replication-shop-benchmark/logs
BENCHMARK_EXIT_AFTER_TIMEOUT=120

SLEEP_SECONDS=65

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

down_docker_compose_files() {
  local files
  files=$(find "$SCRIPT_DIR"/../ -type f -name 'docker-compose*.yml')

  for file in $files; do
    printf "Stopping %s containers\n" "$file"
    docker compose -f "$file" down
  done
  if [ "$(docker volume ls -q | wc -l)" -ge 1 ]; then
    docker volume rm $(docker volume ls -q)
  fi
}

down_docker_compose_files

for branch_name in "${BRANCHES[@]}"; do
  # shellcheck disable=SC2001
  LOG_FILE_PREFIX=$(echo "$branch_name" | sed 's#/#\-#g')
  export LOG_FILE_PREFIX="$LOG_FILE_PREFIX"

  BRANCH="$branch_name"

  printf "${YELLOW}Branch: %s${NC}\n" "$branch_name"

  printf "${YELLOW}Building images${NC}\n"
  docker build "https://github.com/Lexxxzy/replication-shop.git#$branch_name" \
    -t "replication-shop-app-$branch_name"

  printf "${YELLOW}Start services${NC}\n"
  compose_files=(-f docker-compose.app.yml)
  need_sleep=false
  case "$branch_name" in
  *postgresql*)
    compose_files+=(-f docker-compose.postgresql.yml)
    export POSTGRESQL_ENABLED=true
    export CASSANDRA_ENABLED=false
    export REDIS_ENABLED=false
    ;;
  *cassandra*)
    compose_files+=(-f docker-compose.cassandra.yml)
    need_sleep=true
    export POSTGRESQL_ENABLED=true
    export CASSANDRA_ENABLED=true
    export REDIS_ENABLED=false
    ;;
  *redis*)
    compose_files+=(-f docker-compose.redis.yml)
    export POSTGRESQL_ENABLED=true
    export CASSANDRA_ENABLED=true
    export REDIS_ENABLED=true
    ;;
  esac
  docker compose "${compose_files[@]}" up -d
  BRANCH="-$BRANCH" \
    docker compose -f docker-compose.app.yml up -d

  if [ "$need_sleep" == true ]; then
    SLEEP_SECONDS=$((SLEEP_SECONDS * 3))
  fi

  printf "${YELLOW}Wait for services to start${NC}\n"
  sleep $((SLEEP_SECONDS))

  printf "${YELLOW}Start benchmark${NC}\n"
  docker compose -f docker-compose.benchmark.yml build
  EXIT_AFTER_TIMEOUT="$BENCHMARK_EXIT_AFTER_TIMEOUT" \
    docker compose -f docker-compose.benchmark.yml up -d --scale benchmark="$BENCHMARK_INSTANCES_NUM" benchmark

  printf "${YELLOW}Wait for benchmark containers stop${NC}\n"
  BENCHMARK_EXIT_AFTER_TIMEOUT=$((BENCHMARK_EXIT_AFTER_TIMEOUT + 15))
  sleep $((BENCHMARK_EXIT_AFTER_TIMEOUT))

  printf "${YELLOW}Save logs${NC}\n"
  logs_dir="$HOME/$BENCHMARK_OUTPUT_DIR/$logs_dir"
  [ -d "$logs_dir" ] || mkdir -p "$logs_dir"
  docker compose -f docker-compose.benchmark.yml cp benchmark:/app/logs "$logs_dir"
  printf "${GREEN}Saved logs to %s${NC}\n" "$logs_dir"
  docker compose -f docker-compose.benchmark.yml down
done
