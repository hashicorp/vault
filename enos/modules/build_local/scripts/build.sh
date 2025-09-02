#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -eux -o pipefail

# Install yarn so we can build the UI
npm install --global yarn || true

export CGO_ENABLED=0

root_dir="$(git rev-parse --show-toplevel)"
pushd "$root_dir" > /dev/null

if [ -n "$BUILD_UI" ] && [ "$BUILD_UI" = "true" ]; then
  make ci-build-ui
fi

make ci-build

popd > /dev/null

echo "--> Bundling $BIN_PATH/* to $BUNDLE_PATH"
zip -r -j "$BUNDLE_PATH" "$BIN_PATH/"
