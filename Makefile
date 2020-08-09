BASE_PKG := github.com/awanku/awanku
DOCKER_COMPOSE_TEST := docker-compose -f docker-compose.test.yml
DOCKER_COMPOSE_DEV := docker-compose -f docker-compose.dev.yml

build:
	go build -o ./dist/core-api $(BASE_PKG)/cmd/core-api

run:
	go run $(BASE_PKG)/cmd/core-api

run-dev:
	watchexec --watch . --exts go --signal SIGKILL --restart make run

test:
	./scripts/test.sh

update-deps:
	rm -rf dist/*
	go get -u ./...
	go mod tidy

docker-build:
	docker build -f docker/Production.dockerfile -t awanku/core-api:latest .

docker-push:
	docker tag awanku/core-api:latest docker.awanku.id/awanku/core-api:latest
	docker push docker.awanku.id/awanku/core-api:latest

docker-test:
	$(DOCKER_COMPOSE_TEST) up -d
	$(DOCKER_COMPOSE_TEST) run core-api sh /app/core-api/scripts/docker-test.sh
	$(DOCKER_COMPOSE_TEST) stop

docker-test-clean:
	$(DOCKER_COMPOSE_TEST) down

docker-dev-run:
	$(DOCKER_COMPOSE_DEV) build core-api
	$(DOCKER_COMPOSE_DEV) up -d
	$(DOCKER_COMPOSE_DEV) exec core-api sh /app/core-api/database/up.sh
	$(DOCKER_COMPOSE_DEV) up

docker-dev-clean:
	$(DOCKER_COMPOSE_DEV) down --remove-orphans --volumes

docker-dev-psql:
	$(DOCKER_COMPOSE_DEV) exec maindb psql -U awanku awanku

docker-stop:
	$(DOCKER_COMPOSE_DEV) stop

swagger-generate:
	rm -rf ./dist/swagger-core-api
	go run github.com/swaggo/swag/cmd/swag init --dir ./cmd/core-api --output ./dist/swagger-core-api --parseDependency --parseInternal

swagger-docker-build:
	docker build -f docker/ApiDocs.dockerfile -t awanku/core-api-docs:latest .

swagger-docker-push:
	docker tag awanku/core-api-docs:latest docker.awanku.id/awanku/core-api-docs:latest
	docker push docker.awanku.id/awanku/core-api-docs:latest

swagger-docker-run:
	docker run --rm -p 8888:80 awanku/core-api-docs:latest
