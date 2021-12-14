#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export SSH_SERVER_NAME=cto
export STACK_NAME=support
export PROJECT_PATH=/home/docker/cto

cp ./.env ./"${STACK_NAME}"/.env
cd "${STACK_NAME}"

docker-compose -f docker-compose.yml config > "${STACK_NAME}".yml

scp "${STACK_NAME}".yml "${SSH_SERVER_NAME}:${PROJECT_PATH}/"

ssh -tt -o StrictHostKeyChecking=no "${SSH_SERVER_NAME}" \
"cd ${PROJECT_PATH} && docker stack deploy -c ${STACK_NAME}.yml $STACK_NAME"

rm .env
