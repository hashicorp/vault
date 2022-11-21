#!/bin/env bash
set -eu -o pipefail

pushd "$(git rev-parse --show-toplevel)" > /dev/null

version=$(<../../../../.release/VERSION)
echo "$version" | sed 's/ //g'

popd > /dev/null
