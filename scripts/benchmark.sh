#!/usr/bin/env bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

source "$SCRIPT_DIR"/../.env

BENCHMARK_INSTANCES_NUM=${BENCHMARK_INSTANCES_NUM:-6}
BENCHMARK_OUTPUT_DIR=replication-shop-benchmark/logs
BENCHMARK_EXIT_AFTER_TIMEOUT=120

SLEEP_SECONDS=180

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
  BRANCH="-$BRANCH" \
    docker build "https://github.com/Lexxxzy/replication-shop.git#$branch_name" \
    -t "replication-shop-app-$branch_name:latest"

  printf "${YELLOW}Start services${NC}\n"
  compose_files=(-f docker-compose.app.yml)
  need_sleep=false
  case "$branch_name" in
  *postgresql*)
    compose_files+=(-f docker-compose.postgresql.yml)
    ;;
  *cassandra*)
    compose_files+=(-f docker-compose.cassandra.yml)
    need_sleep=true
    ;;
  *redis*)
    compose_files+=(-f docker-compose.redis.yml)
    ;;
  esac
  docker compose "${compose_files[@]}" up -d
  docker compose -f docker-compose.app.yml up -d

  if [ "$need_sleep" == true ]; then
    printf "${YELLOW}Wait for services to start${NC}\n"
    sleep $((SLEEP_SECONDS))
  fi

  printf "${YELLOW}Start benchmark${NC}\n"
  BRANCH="-$BRANCH" \
    docker compose -f docker-compose.benchmark.yml build
  EXIT_AFTER_TIMEOUT="$BENCHMARK_EXIT_AFTER_TIMEOUT" \
    docker compose -f docker-compose.benchmark.yml up -d --scale benchmark="$BENCHMARK_INSTANCES_NUM" benchmark

  printf "${YELLOW}Wait for benchmark containers stop${NC}\n"
  BENCHMARK_EXIT_AFTER_TIMEOUT=$((BENCHMARK_EXIT_AFTER_TIMEOUT + 15))
  sleep $((BENCHMARK_EXIT_AFTER_TIMEOUT))

  printf "${YELLOW}Save logs${NC}\n"
  logs_dir="$HOME/$BENCHMARK_OUTPUT_DIR/$logs_dir"
  [ -d "$logs_dir" ] || mkdir -p "$logs_dir"
  docker compose -f docker-compose.benchmark.yml cp benchmark:/app/logs "$BENCHMARK_OUTPUT_DIR"/"$logs_dir"
  printf "${GREEN}Saved logs to %s${NC}\n" "$logs_dir"
  docker compose -f docker-compose.benchmark.yml down
done
