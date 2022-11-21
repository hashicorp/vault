#!/bin/env bash
set -eu -o pipefail

pushd "$(git rev-parse --show-toplevel)" > /dev/null
make crt-get-version
popd > /dev/null
