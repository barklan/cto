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

function _export_pg {
    . local.env
    export POSTGRES_DB POSTGRES_PASSWORD POSTGRES_USER POSTGRES_HOST
}

function _dc {
    export DOCKER_BUILDKIT=1
    docker-compose ${TTY} "${@}"
}

# ----------------------------------------------------------------------------

function up:c {
    export CTO_DATA_PATH=/home/barklan/dev/cto/.cache
    export CTO_MEDIA_PATH=.cache/media
    export CTO_LOCAL_ENV=true
    export CONFIG_ENV=dev
    _export_pg
    export POSTGRES_HOST=localhost:5432
    echo "$POSTGRES_HOST"
    go run cmd/cto/main.go
}

function up:p {
    _export_pg
    export CONFIG_ENV=dev
    export POSTGRES_HOST=localhost:5432
    echo "$POSTGRES_HOST"
    go run cmd/porter/main.go
}

function up:db {
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile db up --build
}

function psql {
    _dc exec db psql -U postgres -d app "${@}"
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
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile main build --parallel
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile main up
}

function direct {
    . .env
    export REGISTRY_PASSWORD REGISTRY_USERNAME
    bash scripts/directswarm.sh
}

function direct:s {
    . .env
    export REGISTRY_PASSWORD REGISTRY_USERNAME
    bash scripts/directsupport.sh
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
    # docker run -it --rm -p 80:80 \
    # -v "$(pwd)"/docs:/usr/share/nginx/html/swagger/ \
    # -e SPEC_URL=swagger/openapi.yml redocly/redoc:v2.0.0-rc.59
    docker run -p 80:8080 -e SWAGGER_JSON=/docs/openapi.yml -v "$(pwd)"/docs:/docs swaggerapi/swagger-ui
}

function docs:bundle {
    docker run --rm -v "$(pwd)"/docs:/spec redocly/openapi-cli bundle -o bundle.json --ext json openapi.yml
}

function proto {
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/protos/"${1}".proto
}

function db:makemigrations {
    docker run -v "$(pwd)"/db/migrations:/migrations --network host migrate/migrate \
    create -ext sql -dir /migrations -seq "${@}"
}

function db:migrate {
    docker run -v "$(pwd)"/db/migrations:/migrations --network host migrate/migrate \
    -database postgres://postgres:postgres@localhost:5432/app?sslmode=disable -path /migrations "${@}"
}

function db:migrate:remote {
    . .env
    export POSTGRES_PASSWORD
    # FIXME
    ssh -tt -o StrictHostKeyChecking=no cto \
    "docker run --network traefik-public migrate/migrate \
    -source github://barklan/cto/db/migrations \
    -database postgres://postgres:${POSTGRES_PASSWORD}@cto_db:5432/app?sslmode=disable up"
}

# -----------------------------------------------------------------------------

function help {
    printf "%s <task> [args]\n\nTasks:\n" "${0}"

    compgen -A function | grep -v "^_" | cat -n
}

TIMEFORMAT=$'\nTask completed in %3lR'
time "${@:-help}"
