#!/usr/bin/env bash

set -eo pipefail

DC="${DC:-exec}"

# If we're running in CI we need to disable TTY allocation for docker-compose
# commands that enable it by default, such as exec and run.
TTY=""
if [[ ! -t 1 ]]; then
    TTY="-T"
fi

# -----------------------------------------------------------------------------
# Helper functions start with _ and aren't listed in this script's help menu.
# -----------------------------------------------------------------------------

function _dc {
    export DOCKER_BUILDKIT=1
    docker-compose ${TTY} "${@}"
}

# ----------------------------------------------------------------------------
# * General purpose local functions.

function up {
    export CTO_DATA_PATH=/home/barklan/dev/gitlab_workflow_bot/.cache
    export CTO_MEDIA_PATH=.cache/media
    export CTO_LOCAL_ENV=true
    export CONFIG_ENV=dev
    go run cmd/cto/main.go
}

function reset {
    rm -r .cache/main
    rm -r .cache/log
}

function front {
    cd frontend && pnpm dev
}

function upd {
    export DOCKER_BUILDKIT=1
    docker-compose -f docker-compose.yml -f docker-compose.local.yml up --build backend
}

function build {
    docker build -t "barklan/gitlab-workflow-bot:$1" .
}

function push {
    docker image push "barklan/gitlab-workflow-bot:$1"
}

function direct {
    bash scripts/direct.sh
}

function proto {
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/protos/main.proto
}


# ----------------------------------
# * Non-local functions



# -----------------------------------------------------------------------------

function help {
    printf "%s <task> [args]\n\nTasks:\n" "${0}"

    compgen -A function | grep -v "^_" | cat -n
}

TIMEFORMAT=$'\nTask completed in %3lR'
time "${@:-help}"
