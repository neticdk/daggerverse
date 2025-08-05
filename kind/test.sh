#!/usr/bin/env bash

set -eu

dagger develop

pushd tests
dagger develop
dagger call run --socket=/var/run/docker.sock
popd