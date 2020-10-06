### Build container
FROM openjdk:11 AS bridge-{{.Params.name}}-backend-build

ARG ARTIFACTORY_USERNAME
ARG ARTIFACTORY_PASSWORD

ENV ORG_GRADLE_PROJECT_artifactoryUsername=$ARTIFACTORY_USERNAME \
    ORG_GRADLE_PROJECT_artifactoryPassword=$ARTIFACTORY_PASSWORD

RUN useradd -u 9999 -l docker
RUN mkdir -p /gradle && chown docker:docker /gradle

USER docker

VOLUME /gradle

COPY --chown=docker . /usr/src/app
WORKDIR /usr/src/app

RUN ./gradlew --no-daemon --gradle-user-home=/gradle --info assemble


### Application container
FROM openjdk:11 AS bridge-{{.Params.name}}-backend-app

COPY --from=bridge-{{.Params.name}}-backend-build --chown=root \
		/usr/src/app/build/libs/bridge-{{.Params.name}}-backend-*-SNAPSHOT.jar \
		/app/bridge-{{.Params.name}}-backend.jar

RUN useradd -u 9999 -l docker
USER docker

COPY ./dev/docker/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

COPY ./dev/docker/start.sh /start.sh
CMD ["/start.sh"]
