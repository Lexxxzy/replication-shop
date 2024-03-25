#!/usr/bin/env bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

docker compose -f "$SCRIPT_DIR"/../docker-compose.jupyter.yml run --rm --build jupyter-nbconvert