#!/usr/bin/env bash


LOG_DIR=${LOG_DIR:-./"$LOG_DIR_PREFIX"logs/}

rm -rf "$LOG_DIR"

mkdir -p "$LOG_DIR"

dotnet run ShopClient.dll --config nginx_config.prod.json --log "$LOG_DIR"/load_test_"$(hostname -f)".csv