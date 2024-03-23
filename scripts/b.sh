#!/usr/bin/env bash

docker build -f ./benchmark/Dockerfile https://github.com/Lexxxzy/replication-shop.git#"${BRANCH:-dev}" -t replication-shop-benchmark:latest
