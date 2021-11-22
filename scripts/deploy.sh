#!/bin/bash

set -e

export ACCESS_KEY=${ACCESS_KEY?Variable not set}
export IP_ADDRESS=${IP_ADDRESS?Variable not set}

eval $(ssh-agent -s)
echo "$ACCESS_KEY" | tr -d '\r' | ssh-add -
mkdir -p ~/.ssh
chmod 700 ~/.ssh
echo $IP_ADDRESS
ssh-keyscan $IP_ADDRESS >> ~/.ssh/known_hosts
chmod 644 ~/.ssh/known_hosts
ssh-keyscan -H 'gitlab.com' >> ~/.ssh/known_hosts

ssh -tt -o StrictHostKeyChecking=no "root@$IP_ADDRESS" \
"mkdir -p /home/docker/gitlab/environment && cd /home/docker/gitlab && docker-compose down"
scp docker-compose.yml "root@$IP_ADDRESS:/home/docker/gitlab/"
scp environment/cto_nft.yml "root@$IP_ADDRESS:/home/docker/gitlab/environment/cto.yml"

ssh -tt -o StrictHostKeyChecking=no "root@$IP_ADDRESS" \
"cd /home/docker/gitlab && \
mkdir -p .cache/media && \
docker volume create cto-data && \
docker image rm barklan/gitlab-workflow-bot:rolling || true && \
docker-compose up -d"
