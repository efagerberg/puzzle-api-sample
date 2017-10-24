BASH ?= /bin/bash

DC_COMMAND = docker-compose
DC_UP_COMMAND = $(DC_COMMAND) up
DC_RUN_COMMAND = $(DC_COMMAND) run --rm

## General targets
default:
	export APP_DB_USER=root
	export APP_DB_NAME=puzzle_api
	export APP_DB_HOST=192.168.99.100
	export APP_DB_PORT=5432
	postgres/provision-database.sh
	echo "$(APP_DB_HOST)"
	postgres/wait-for-postgres.sh $(APP_DB_HOST)
	go run main.go

tests:
	postgres/provision-database.sh
	postgres/wait-for-postgres.sh $(TEST_DB_HOST)
	go test

nuke:
	docker rm -f $$(docker ps -aqf "name=puzzle-api-sample") || echo $$'No containers to remove.'

## General docker targets
docker-bash:
	$(DC_RUN) app $(BASH)

docker-stop:
	$(DC_CMD) stop

docker-build:
	$(DC_CMD) build

docker-logs:
	$(DC_CMD) logs -f --tail=30