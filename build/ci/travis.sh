#!/bin/bash

set -e

# Do not rebuild/retest image that we already have.
if [[ -n "$TRAVIS_TAG" ]]; then
  echo "Found release tag"
  docker pull "${DOCKER_USERNAME}/janna:${COMMIT}"
  docker tag "${DOCKER_USERNAME}/janna:${COMMIT}" "${DOCKER_USERNAME}/janna:${TRAVIS_TAG}"
  docker push "${DOCKER_USERNAME}/janna:${TRAVIS_TAG}"
  exit 0
fi

make

if [[ "$TRAVIS_BRANCH" == "master" ]]; then
  make push
  make push TAG=latest
fi

docker save "${DOCKER_USERNAME}/janna:${TRAVIS_TAG}" | gzip -c > "${CACHE_FILE}"
