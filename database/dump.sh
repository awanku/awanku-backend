#!/usr/bin/env sh
if [[ "$DATABASE_URL" == '' ]]; then
    echo 'please set database url as DATABASE_URL in environment variable'
    exit 1
fi

pg_dump --file=database/dump.sql --verbose "${DATABASE_URL}"
