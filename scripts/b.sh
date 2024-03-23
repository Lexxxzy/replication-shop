#!/usr/bin/env bash

BRANCHES=("feature/cassandra-cluster" "feature/postgresql-cluster" "feature/redis")

for branch_name in "${BRANCHES[@]}"; do
    printf "Building image for branch: %s\n" "$branch_name"
    docker build -f ./benchmark/Dockerfile "https://github.com/Lexxxzy/replication-shop.git#$branch_name"\
     -t "replication-shop-benchmark-$branch_name:latest"
done
