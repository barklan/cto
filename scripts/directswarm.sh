#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export SSH_SERVER_NAME=cto
export STACK_NAME=cto
export PROJECT_PATH=/home/docker/cto

docker build -t "barklan/cto-core:rolling" .
docker image push "barklan/cto-core:rolling"
docker build -t "barklan/cto-explorer:rolling" ./frontend
docker image push "barklan/cto-explorer:rolling"

docker-compose -f docker-compose.yml config > docker-stack.yml

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" "mkdir -p ${PROJECT_PATH}/environment"
scp .env "${SSH_SERVER_NAME}:${PROJECT_PATH}"

scp docker-stack.yml "${SSH_SERVER_NAME}:${PROJECT_PATH}/"
scp -r environment "${SSH_SERVER_NAME}:${PROJECT_PATH}"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"cd ${PROJECT_PATH} && docker stack deploy -c docker-stack.yml $STACK_NAME"
