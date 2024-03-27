#!/usr/bin/env bash

set -e

cqlsh --request-timeout=6000 -f /statements/ddl/cassandra.cql
cqlsh --request-timeout=6000 -f /statements/dml/cassandra.cql