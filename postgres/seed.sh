#!/usr/bin/env bash
echo "***** Creating hstore extension *****"
gosu postgres psql --dbname template1 <<- EOSQL
    CREATE EXTENSION hstore;
    DROP DATABASE postgres;
    CREATE DATABASE postgres TEMPLATE template1;
EOSQL
echo "***** Creating puzzle DB *****"
gosu postgres psql <<- EOSQL
    CREATE USER evan;
    ALTER USER evan CREATEDB;

    CREATE DATABASE puzzle;
    GRANT ALL PRIVILEGES ON DATABASE puzzle TO evan;
EOSQL
