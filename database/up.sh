#!/usr/bin/env sh
if [[ "$DATABASE_URL" == '' ]]; then
    echo 'please set database url as DATABASE_URL in environment variable'
    exit 1
fi
go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -database "${DATABASE_URL}" -source file://./database/migrations up
