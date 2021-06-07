#!/usr/bin/env bash

set -xeou pipefail

cd /opt/bonvoy || exit 1

export DOCKER_API_VERSION=1.39

go mod tidy
go build bonvoy