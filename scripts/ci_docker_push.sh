#!/bin/bash -eu
set -x

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push "orgplace/trafficquota:${TRAVIS_TAG:-latest}"
