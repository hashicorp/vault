#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${KMIP_PORT}" ]] && fail "KMIP_PORT env variable has not been set"

# Pull KMIP Docker image
CONTAINER_CMD="sudo podman"
KMIP_DOCKER_NAME="docker.io/percona/percona-server:8.0"
${CONTAINER_CMD} pull "${KMIP_DOCKER_NAME}"

mkdir TEMP_DIR
cd TEMP_DIR
TEMP_DIR=$(pwd)

# Run KMIP container
echo "Starting KMIP container..."
${CONTAINER_CMD} run -d \
  --name kmip \
  --volume "${TEMP_DIR}":/TEMP_DIR \
  -e KMIP_ADDR="${VAULT_ADDR}" \
  -e MYSQL_ROOT_PASSWORD=testpassword \
  "${KMIP_DOCKER_NAME}" \
  --port "${KMIP_PORT}"

echo "KMIP server is now running in Docker!"
