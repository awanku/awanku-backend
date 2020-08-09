#!/usr/bin/env sh
if [[ "$DATABASE_URL" == '' ]]; then
    echo 'please set database url as DATABASE_URL in environment variable'
    exit 1
fi

psql "${DATABASE_URL}" <<-EOF
drop schema public cascade;
create schema public;
EOF
