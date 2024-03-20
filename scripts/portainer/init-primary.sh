#!/usr/bin/env bash

docker swarm init
docker network create -d overlay --attachable public
curl -L https://downloads.portainer.io/ce2-19/portainer-agent-stack.yml -o portainer-agent-stack.yml
docker stack deploy -c portainer-agent-stack.yml portainer