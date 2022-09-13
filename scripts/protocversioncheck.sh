#!/usr/bin/env bash

set -euo pipefail

PROTOC_CMD=${PROTOC_CMD:-protoc}
PROTOC_VERSION_MIN="$1"
echo "==> Checking that protoc is at version $1..."

PROTOC_VERSION=$($PROTOC_CMD --version | grep -o '[0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?')

IFS="." read -r -a PROTOC_VERSION_ARR <<< "$PROTOC_VERSION"
IFS="." read -r -a PROTOC_VERSION_REQ <<< "$PROTOC_VERSION_MIN"

if [[ ${PROTOC_VERSION_ARR[0]} -lt ${PROTOC_VERSION_REQ[0]} ||
  ( ${PROTOC_VERSION_ARR[0]} -eq ${PROTOC_VERSION_REQ[0]} &&
  ( ${PROTOC_VERSION_ARR[1]} -lt ${PROTOC_VERSION_REQ[1]} ||
  ( ${PROTOC_VERSION_ARR[1]} -eq ${PROTOC_VERSION_REQ[1]} && ${PROTOC_VERSION_ARR[2]} -lt ${PROTOC_VERSION_REQ[2]} )))
  ]]; then
  echo "protoc should be at $PROTOC_VERSION_MIN; found $PROTOC_VERSION."
  exit 1
fi

echo "Using protoc version $PROTOC_VERSION"
