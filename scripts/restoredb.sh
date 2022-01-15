#!/bin/bash

docker exec "$(docker ps -q -f name=cto_db)" psql -U postgres -c "DROP DATABASE app WITH (FORCE);"
docker exec "$(docker ps -q -f name=cto_db)" bash -c "createdb -U postgres -T template0 app"
docker exec -i "$(docker ps -q -f name=cto_db)" psql -U postgres app < ./db_dump.sql
docker exec "$(docker ps -q -f name=cto_db)" psql -U postgres -c "GRANT CONNECT ON DATABASE app TO public;"
