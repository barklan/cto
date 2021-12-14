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

function up {
    export CTO_DATA_PATH=/home/barklan/dev/cto/.cache
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
    docker-compose -f docker-compose.yml -f docker-compose.local.yml up --build
}

function direct {
    bash scripts/directswarm.sh
}

function fluentd:push {
    cd dockerfiles/fluentd_cto
    docker build -t barklan/fluentd-cto:"$1" .
    docker image push barklan/fluentd-cto:"$1"

    cd ../fluentd_cto_es
    docker build -t barklan/fluentd-cto:"$1"es .
    docker image push barklan/fluentd-cto:"$1"es
}

function docs:dev {
    docker run -it --rm -p 80:80 \
    -v "$(pwd)"/docs:/usr/share/nginx/html/swagger/ \
    -e SPEC_URL=swagger/openapi.yml redocly/redoc:v2.0.0-rc.59
}

function docs:bundle {
    cd docs && npx redoc-cli bundle openapi.yml && mv redoc-static.html index.html
}

# -----------------------------------------------------------------------------

function help {
    printf "%s <task> [args]\n\nTasks:\n" "${0}"

    compgen -A function | grep -v "^_" | cat -n
}

TIMEFORMAT=$'\nTask completed in %3lR'
time "${@:-help}"
