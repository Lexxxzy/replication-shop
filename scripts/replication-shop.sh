#!/usr/bin/env bash

# Default values
SVC_CMD=""
APP_CMD=""

show_help() {
    echo "Usage: $0 -a [action]"
    echo "Options:"
    echo "  -a, --action   Specify the action to perform. Valid values: 'up' or 'down'."
    echo "  -h, --help     Show this help message and exit."
}

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -a|--action)
            ACTION="$2"
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Error: Unknown option $1"
            show_help
            exit 1
            ;;
    esac
    shift
done

SVC_UP_CMD="up -d"
APP_UP_CMD="up -d --build"
SVC_DOWN_CMD="down"
APP_DOWN_CMD="down"

if [ "$ACTION" == "up" ]; then
    SVC_CMD="$SVC_UP_CMD"
    APP_CMD="$APP_UP_CMD"
elif [ "$ACTION" == "down" ]; then
    SVC_CMD="$SVC_DOWN_CMD"
    APP_CMD="$APP_DOWN_CMD"
fi

docker compose \
    -f docker-compose.postgresql.yml \
    -f docker-compose.redis.yml \
    -f docker-compose.cassandra.yml \
    $SVC_CMD

docker compose \
    -f docker-compose.app.yml \
    $APP_CMD

