#!/bin/bash

set -e
start=$(date +%s)

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

# save layers for cache
docker save "${DOCKER_USERNAME}/janna:stage-env" --output="${CACHE_FILE_STAGE_ENV}"
docker save "${DOCKER_USERNAME}/janna:latest" --output="${CACHE_FILE}"

end=$(date +%s)

runtime=$((end-start))
echo "$0 took: ${runtime}s"
