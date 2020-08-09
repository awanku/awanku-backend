#!/usr/bin/env sh
name="${1}"
if [[ "$name" == '' ]]; then
    echo 'please specify migration name'
    exit 1
fi
go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir database/migrations "${name}"
