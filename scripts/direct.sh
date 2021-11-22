#!/bin/bash

set -e

export DOCKER_BUILDKIT=1

docker build -t "barklan/gitlab-workflow-bot:rolling" .
docker image push "barklan/gitlab-workflow-bot:rolling"
docker build -t "barklan/cto-explorer:rolling" ./frontend
docker image push "barklan/cto-explorer:rolling"

ssh -tt -o StrictHostKeyChecking=no "helper" "mkdir -p /home/docker/cto/.cache && mkdir -p /home/docker/cto/environment"
scp -r .cache/media "helper:/home/docker/cto/.cache"
scp .env "helper:/home/docker/cto"

ssh -tt -o StrictHostKeyChecking=no "helper" \
"mkdir -p /home/docker/cto && cd /home/docker/cto && docker-compose down"
scp docker-compose.yml "helper:/home/docker/cto/"
scp -r environment "helper:/home/docker/cto"

ssh -tt -o StrictHostKeyChecking=no "helper" \
"cd /home/docker/cto && \
mkdir -p .cache/media && \
docker volume create cto-data && \
docker image rm barklan/gitlab-workflow-bot:rolling || true && \
docker image rm barklan/cto-explorer:rolling || true && \
docker-compose up -d"
