### Build container
FROM instructure/java:11 AS bridge-{{.Params.name}}-backend-build

ARG ARTIFACTORY_USERNAME
ARG ARTIFACTORY_PASSWORD

ENV ORG_GRADLE_PROJECT_artifactoryUsername=$ARTIFACTORY_USERNAME \
	ORG_GRADLE_PROJECT_artifactoryPassword=$ARTIFACTORY_PASSWORD \
	GRADLE_OPTS="-Dorg.gradle.internal.launcher.welcomeMessageEnabled=false"

RUN mkdir -p ~/.gradle && chown docker:docker ~/.gradle

COPY --chown=docker . /usr/src/app
WORKDIR /usr/src/app

RUN ./gradlew --no-daemon assemble

### Appliaction container
FROM instructure/java:11 AS bridge-{{.Params.name}}-backend

COPY --from=bridge-{{.Params.name}}-backend-build \
		/usr/src/app/build/libs/bridge-{{.Params.name}}-backend-*-SNAPSHOT.jar \
		/app/bridge-{{.Params.name}}-backend.jar

USER docker

COPY ["scripts/entrypoint.sh", "scripts/start.sh", "/app/"]
WORKDIR /app

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/start.sh"]
