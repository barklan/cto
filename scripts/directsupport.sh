#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export SSH_SERVER_NAME=cto
export STACK_NAME=support
export PROJECT_PATH=/home/docker/cto
export REGISTRY_USERNAME="${REGISTRY_USERNAME?Variable not set}"
export REGISTRY_PASSWORD="${REGISTRY_PASSWORD?Variable not set}"
export DOCKER_REGISTRY="${DOCKER_REGISTRY?Variable not set}"
export DOCKER_IMAGE_PREFIX="${DOCKER_IMAGE_PREFIX?Variable not set}"

docker login -u "${REGISTRY_USERNAME}" -p "${REGISTRY_PASSWORD}" "${DOCKER_REGISTRY}"

cp ./.env ./"${STACK_NAME}"/.env
cd "${STACK_NAME}"

(echo -e "version: '3.9'\n";  docker compose -f docker-compose.yml config) > "${STACK_NAME}".yml

scp "${STACK_NAME}".yml "${SSH_SERVER_NAME}:${PROJECT_PATH}/"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"docker login -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD} ${DOCKER_REGISTRY} \
&& cd ${PROJECT_PATH} && docker stack deploy -c ${STACK_NAME}.yml --with-registry-auth $STACK_NAME"

rm .env
