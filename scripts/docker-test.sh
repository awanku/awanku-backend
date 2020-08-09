#!/usr/bin/env sh
name='awanku-test-1'
psql "${OWNER_DB_URL}" <<SQL
drop database if exists "${name}";
drop role if exists "${name}";
create role "${name}" with login password 'rahasia';
create database "${name}" owner '${name}';
SQL

export DATABASE_URL="postgres://${name}:rahasia@${DB_HOST}/${name}?sslmode=disable"
./database/up.sh

go test -v -race -cover $(go list ./... | grep -v dist)
