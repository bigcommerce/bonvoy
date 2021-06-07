#!/usr/bin/env bash

TARGET=${1:-"auth-grpc"}

cd /opt/bonvoy || exit 1

export DOCKER_API_VERSION=1.39

go mod tidy
go build bonvoy
./bonvoy listeners "$TARGET"