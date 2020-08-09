#!/usr/bin/env bash
docker tag awanku/core-api:latest docker.awanku.id/awanku/core-api:latest
docker tag awanku/core-api:latest docker.awanku.id/awanku/core-api:${DOCKER_IMAGE_TAG}

docker push docker.awanku.id/awanku/core-api:latest
docker push docker.awanku.id/awanku/core-api:${DOCKER_IMAGE_TAG}
