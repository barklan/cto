#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export SSH_SERVER_NAME=cto

docker build -t "barklan/cto-core:rolling" .
docker image push "barklan/cto-core:rolling"
docker build -t "barklan/cto-explorer:rolling" ./frontend
docker image push "barklan/cto-explorer:rolling"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" "mkdir -p /home/docker/cto/.cache && mkdir -p /home/docker/cto/environment"
scp .env "${SSH_SERVER_NAME}:/home/docker/cto"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"mkdir -p /home/docker/cto && cd /home/docker/cto && docker-compose down"
scp docker-compose.yml "${SSH_SERVER_NAME}:/home/docker/cto/"
scp -r environment "${SSH_SERVER_NAME}:/home/docker/cto"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"cd /home/docker/cto && \
docker volume create cto-data && \
docker volume create cto-media && \
docker image rm barklan/cto-core:rolling || true && \
docker image rm barklan/cto-explorer:rolling || true && \
docker-compose up -d"
