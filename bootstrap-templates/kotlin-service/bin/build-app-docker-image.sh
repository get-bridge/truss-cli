#!/bin/sh

set -eu

: ${BASEDIR:="$(git rev-parse --show-toplevel)"}
: ${ARTIFACTORY_USERNAME:="$($BASEDIR/gradlew -q printProperty -PpropertyName=artifactoryUsername)"}
: ${ARTIFACTORY_PASSWORD:="$($BASEDIR/gradlew -q printProperty -PpropertyName=artifactoryPassword)"}

docker build \
	--force-rm \
	--build-arg ARTIFACTORY_USERNAME="$ARTIFACTORY_USERNAME" \
	--build-arg ARTIFACTORY_PASSWORD="$ARTIFACTORY_PASSWORD" \
	-t bridge-{{.Params.name}}-backend \
	-f "${BASEDIR}/dev/docker/Dockerfile" \
	"${BASEDIR}"
