#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${SERVER_CA}" ]] && fail "SERVER_CA env variable has not been set"
[[ -z "${CLIENT_CA}" ]] && fail "CLIENT_CA env variable has not been set"
[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${KMIP_PORT}" ]] && fail "KMIP_PORT env variable has not been set"

cd ~ || fail "Failed to change directory to home"
echo "${SERVER_CA}" > TEMP_DIR/vault-ca.pem
echo "${CLIENT_CA}" > TEMP_DIR/client.pem

# Extract certificate and key from client bundle
cd TEMP_DIR
# Assuming CLIENT_CA contains both cert and key, split them
csplit -f client- client.pem '/-----BEGIN.*PRIVATE KEY-----/' '{*}'
mv client-00 cert.pem
mv client-01 key.pem

# Connect to the Percona Docker container
CONTAINER_CMD="sudo docker"
KMIP_DOCKER_NAME="kmip"

# Create MySQL data directory
${CONTAINER_CMD} exec -d "${KMIP_DOCKER_NAME}" sh -c 'mkdir -p /var/lib/mysql/testKMIP'

# Start MySQL with KMIP configuration
${CONTAINER_CMD} exec -d "${KMIP_DOCKER_NAME}" mysqld \
    --datadir=/var/lib/mysql/testKMIP \
    --early-plugin-load=keyring_kmip.so \
    --keyring_kmip_server_name="${VAULT_ADDR}" \
    --keyring_kmip_server_port="${KMIP_PORT}" \
    --keyring_kmip_client_ca=/TEMP_DIR/vault-ca.pem \
    --keyring_kmip_client_key=/TEMP_DIR/key.pem \
    --keyring_kmip_client_cert=/TEMP_DIR/cert.pem
