#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# This script builds Vault binaries and optionally packages them.
#
# Two distinct workflows are supported:
# 1. Standard build: Builds to dist/ and creates zip bundle from dist/
# 2. Target path build: Builds to dist/, copies to TARGET_BIN_PATH, skips bundling
#    (bundling is skipped to avoid confusion when binary exists in multiple locations)
#
# Environment variables:
# - BUILD_UI: Set to "true" to build UI components
# - TARGET_BIN_PATH: If set, copies built binary to this location instead of bundling
# - BUNDLE_PATH: If set (and TARGET_BIN_PATH is not), creates zip bundle at this path

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

if [ -n "$TARGET_BIN_PATH" ]; then
  echo "--> Target binary path specified, copying binary and skipping bundle"
  make -C "$root_dir" ci-copy-binary
elif [ -n "$BUNDLE_PATH" ]; then
  echo "--> Creating zip bundle from dist/"
  make -C "$root_dir" ci-bundle
else
  echo "--> No post-build packaging requested (neither TARGET_BIN_PATH nor BUNDLE_PATH specified)"
fi
