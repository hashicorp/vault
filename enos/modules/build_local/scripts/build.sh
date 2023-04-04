#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

set -eux -o pipefail

# Install yarn so we can build the UI
npm install --global yarn || true

export CGO_ENABLED=0

root_dir="$(git rev-parse --show-toplevel)"
pushd "$root_dir" > /dev/null
make ci-build-ui ci-build ci-bundle
popd > /dev/null
