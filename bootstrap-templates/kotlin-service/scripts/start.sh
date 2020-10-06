#!/bin/sh

exec java \
    -XX:+UseContainerSupport \
    -XX:MaxRAMPercentage=80 \
    -jar /app/bridge-{{.Params.name}}-backend.jar
