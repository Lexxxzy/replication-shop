#!/usr/bin/env bash

LOG_FILE_PREFIX=$(git branch --show-current)
SCALE_BENCHMARK=${SCALE_BENCHMARK:-6}

LOG_FILE_PREFIX=$LOG_FILE_PREFIX docker compose -f docker-compose.benchmark.yml up \
-d --build \
--scale benchmark="$SCALE_BENCHMARK" \
benchmark