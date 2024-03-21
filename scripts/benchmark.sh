#!/usr/bin/env bash

SCALE_BENCHMARK=${SCALE_BENCHMARK:-2}

docker compose -f docker-compose.benchmark.yml up \
-d --build \
--scale benchmark="$SCALE_BENCHMARK" \
benchmark