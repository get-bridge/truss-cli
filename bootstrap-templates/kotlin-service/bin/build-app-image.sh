#!/bin/sh

set -euo pipefail

REPOSITORY=127178877223.dkr.ecr.us-east-2.amazonaws.com/bridge-{{.Params.name}}/backend

: ${BASEDIR:="$(cd "$(dirname ${BASH_SOURCE[0]:-$0})/.."; pwd)"}
cd "$BASEDIR"

: ${ARTIFACTORY_USERNAME:="$($BASEDIR/gradlew -q printProperty -PpropertyName=artifactoryUsername)"}
: ${ARTIFACTORY_PASSWORD:="$($BASEDIR/gradlew -q printProperty -PpropertyName=artifactoryPassword)"}

GIT_COMMIT=$(cd "$BASEDIR" && [[ "$(git rev-parse HEAD 2>/dev/null)" != "HEAD" ]] && git rev-parse HEAD || true)

docker build \
	--force-rm \
	--build-arg ARTIFACTORY_USERNAME="$ARTIFACTORY_USERNAME" \
	--build-arg ARTIFACTORY_PASSWORD="$ARTIFACTORY_PASSWORD" \
	-t "$REPOSITORY:latest" \
	-f "${BASEDIR}/Dockerfile" \
	"${BASEDIR}"

docker push "$REPOSITORY:latest"

if [[ -n "$GIT_COMMIT" ]]; then
    docker tag "$REPOSITORY:latest" "$REPOSITORY:$GIT_COMMIT"
    docker push "$REPOSITORY:$GIT_COMMIT"
fi
