#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -euo pipefail

PROTOC_CMD=${PROTOC_CMD:-protoc}
PROTOC_VERSION_EXACT="$1"
echo "==> Checking that protoc is at version $1..."

PROTOC_VERSION=$($PROTOC_CMD --version | grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+')

if [ "$PROTOC_VERSION" == "$PROTOC_VERSION_EXACT" ]; then
  echo "Using protoc version $PROTOC_VERSION"
else
  echo "protoc should be at $PROTOC_VERSION_EXACT; found $PROTOC_VERSION."
  echo "If your version is higher than the version this script is looking for, updating the Makefile with the newer version."
  exit 1
fi
