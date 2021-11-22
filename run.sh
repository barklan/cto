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

function run {
    export CTO_DATA_PATH=/home/barklan/dev/gitlab_workflow_bot/.cache
    export CTO_MEDIA_PATH=.cache/media
    export CTO_LOCAL_ENV=true
    go run main.go
}

function reset {
    rm -r .cache/main
    rm -r .cache/log
}

function front {
    cd frontend && pnpm dev
}

function runD {
    export DOCKER_BUILDKIT=1
    docker-compose -f docker-compose.local.yml up --build
}

function build {
    docker build -t "barklan/gitlab-workflow-bot:$1" .
}

function push {
    docker image push "barklan/gitlab-workflow-bot:$1"
}

# DO NOT USE!
function nft {
    bash  scripts/direct.sh
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
