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

function _export_common {
    . .env
    export POSTGRES_DB POSTGRES_PASSWORD POSTGRES_USER
    export RABBITMQ_DEFAULT_USER RABBITMQ_DEFAULT_PASS
    export REDIS_PASSWORD
    export POSTGRES_HOST=localhost:5432
    export RABBITMQ_HOST=localhost
    export REDIS_HOST=localhost
}

function _dc {
    export DOCKER_BUILDKIT=1
    docker-compose ${TTY} "${@}"
}

# ----------------------------------------------------------------------------

up() {
    _export_common
    reflex -c reflex.conf --decoration=fancy
}

up:core() {
    export CTO_DATA_PATH=/home/barklan/dev/cto/.cache
    export CTO_MEDIA_PATH=.cache/media
    export CTO_LOCAL_ENV=true
    export CONFIG_ENV=dev
    _export_common
    go run cmd/cto/main.go
}

up:porter() {
    _export_common
    export CONFIG_ENV=dev
    export OAUTH_CLIENT_ID OAUTH_CLIENT_SECRET
    export OAUTH_CALLBACK_URI=http://localhost:9010/api/porter/signin/callback
    go run cmd/porter/main.go
}

up:loginput() {
    _export_common
    go run cmd/loginput/main.go
}

up:db() {
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile db up --build
}

up:extra() {
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile mq --profile db --profile cache up --build
}

psql() {
    _dc exec db psql -U postgres -d app "${@}"
}

badger:reset() {
    rm -r .cache/main
    rm -r .cache/log
}

frontend() {
    cd frontend && pnpm dev
}

up:fullstack() {
    export DOCKER_BUILDKIT=1
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile main build --parallel
    docker-compose -f docker-compose.yml -f docker-compose.local.yml --profile main up
}

ci:direct() {
    . .env
    export REGISTRY_PASSWORD REGISTRY_USERNAME DOCKER_IMAGE_PREFIX DOCKER_REGISTRY
    echo "$DOCKER_REGISTRY"
    bash scripts/directswarm.sh
}

ci:direct:support() {
    . .env
    export REGISTRY_PASSWORD REGISTRY_USERNAME DOCKER_IMAGE_PREFIX DOCKER_REGISTRY
    bash scripts/directsupport.sh
}

fluentd:push() {
    cd dockerfiles/fluentd_cto
    docker build -t barklan/fluentd-cto:"$1" .
    docker image push barklan/fluentd-cto:"$1"

    cd ../fluentd_cto_es
    docker build -t barklan/fluentd-cto:"$1"es .
    docker image push barklan/fluentd-cto:"$1"es
}

docs:dev() {
    # docker run -it --rm -p 80:80 \
    # -v "$(pwd)"/docs:/usr/share/nginx/html/swagger/ \
    # -e SPEC_URL=swagger/openapi.yml redocly/redoc:v2.0.0-rc.59
    docker run -p 80:8080 -e SWAGGER_JSON=/docs/openapi.yml -v "$(pwd)"/docs:/docs swaggerapi/swagger-ui
}

docs:bundle() {
    docker run --rm -v "$(pwd)"/docs:/spec redocly/openapi-cli bundle -o bundle.json --ext json openapi.yml
}

proto() {
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/protos/"${1}".proto
}

db:makemigrations() {
    docker run -v "$(pwd)"/db/migrations:/migrations --network host migrate/migrate \
    create -ext sql -dir /migrations -seq "${@}"
}

db:migrate() {
    _export_common
    docker run -v "$(pwd)"/db/migrations:/migrations --network host migrate/migrate \
    -database postgres://postgres:${POSTGRES_PASSWORD}@localhost:5432/app?sslmode=disable -path /migrations up
}

db:migrate:remote() {
    . .env
    export POSTGRES_PASSWORD
    ssh -tt -o StrictHostKeyChecking=no cto \
    "docker run --network traefik-public migrate/migrate \
    -source github://barklan/cto/db/migrations \
    -database postgres://postgres:${POSTGRES_PASSWORD}@cto_db:5432/app?sslmode=disable up"
}

log() {
    ssh -tt cto "docker service logs cto_$1 --since $2m"
}

# -----------------------------------------------------------------------------

function help {
    printf "%s <task> [args]\n\nTasks:\n" "${0}"

    compgen -A function | grep -v "^_" | cat -n
}

TIMEFORMAT=$'\nTask completed in %3lR'
time "${@:-help}"
