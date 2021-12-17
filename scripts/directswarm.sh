#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export SSH_SERVER_NAME=cto
export STACK_NAME=cto
export PROJECT_PATH=/home/docker/cto
export REGISTRY_USERNAME="${REGISTRY_USERNAME?Variable not set}"
export REGISTRY_PASSWORD="${REGISTRY_PASSWORD?Variable not set}"

docker build -t "barklan/cto-core:rolling" -f dockerfiles/core.dockerfile .
docker image push "barklan/cto-core:rolling"
docker build -t "barklan/cto-explorer:rolling" -f dockerfiles/frontend.dockerfile ./frontend
docker image push "barklan/cto-explorer:rolling"
docker build -t "barklan/cto-porter:rolling" -f dockerfiles/porter.dockerfile .
docker image push "barklan/cto-porter:rolling"

docker-compose -f docker-compose.yml config > docker-stack.yml

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" "mkdir -p ${PROJECT_PATH}/environment"

scp docker-stack.yml "${SSH_SERVER_NAME}:${PROJECT_PATH}/"
scp -r environment "${SSH_SERVER_NAME}:${PROJECT_PATH}"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"docker login -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD} \
&& cd ${PROJECT_PATH} && docker stack deploy -c docker-stack.yml --with-registry-auth $STACK_NAME"
