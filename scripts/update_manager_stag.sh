#! /usr/bin/env bash

# Exit in case of error
set -e

bash -c "cd ./cmd/manager && go build -o manager ."

scp ./cmd/manager/manager stag:/home/ubuntu/stag/manager
