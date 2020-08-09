#!/usr/bin/env bash
if [[ "${DB_URL}" == '' ]]; then
    export DB_URL='postgres://postgres:@localhost/awanku?sslmode=disable'
fi

./database/nuke.sh
./database/up.sh

go test -v -race -cover $(go list ./... | grep -v dist)
