#! /usr/bin/env bash

# Exit in case of error
set -e


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o ./cmd/manager/manager ./cmd/manager/.

scp ./cmd/manager/manager stag:/home/ubuntu/stag/manager
