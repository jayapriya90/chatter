#!/bin/bash
set -e #exit on error
set -x

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<'EOSQL'
    CREATE DATABASE chatter;
    CREATE USER chatter WITH PASSWORD 'chatter';
    GRANT ALL PRIVILEGES ON DATABASE chatter to chatter;
EOSQL



