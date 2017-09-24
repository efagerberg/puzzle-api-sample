#!/usr/bin/env bash
echo "***** Creating hstore extension *****"
gosu postgres psql --dbname template1 <<- EOSQL
    CREATE EXTENSION hstore;
    DROP DATABASE postgres;
    CREATE DATABASE postgres TEMPLATE template1;
EOSQL
echo "***** Creating puzzle DB *****"
gosu postgres psql <<- EOSQL
    CREATE USER root;
    ALTER USER root CREATEDB;

    CREATE DATABASE puzzle_api;
    GRANT ALL PRIVILEGES ON DATABASE puzzle_api TO root;
EOSQL
