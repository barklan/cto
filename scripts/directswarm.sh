#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export SSH_SERVER_NAME=cto
export STACK_NAME=cto
export PROJECT_PATH=/home/docker/cto
export REGISTRY_USERNAME="${REGISTRY_USERNAME?Variable not set}"
export REGISTRY_PASSWORD="${REGISTRY_PASSWORD?Variable not set}"
export DOCKER_REGISTRY="${DOCKER_REGISTRY?Variable not set}"
export DOCKER_IMAGE_PREFIX="${DOCKER_IMAGE_PREFIX?Variable not set}"

docker login -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD} ${DOCKER_REGISTRY}

docker-compose build #--parallel
docker-compose push

docker-compose -f docker-compose.yml config > docker-stack.yml

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" "mkdir -p ${PROJECT_PATH}/environment"

scp docker-stack.yml "${SSH_SERVER_NAME}:${PROJECT_PATH}/"
scp -r environment "${SSH_SERVER_NAME}:${PROJECT_PATH}"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"docker login -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD} ${DOCKER_REGISTRY} \
&& cd ${PROJECT_PATH} && docker stack deploy -c docker-stack.yml --with-registry-auth $STACK_NAME"
