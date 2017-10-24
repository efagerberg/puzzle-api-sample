#!/bin/bash
docker build -t "puzzle-database:latest" postgres/
if [ ! "$(docker ps -q -f name=puzzle-database)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=puzzle-database)" ]; then
        # startup
        docker start puzzle-database
    fi
    # run your container
    docker run -p 5432:5432 -d --name puzzle-database puzzle-database:latest
fi