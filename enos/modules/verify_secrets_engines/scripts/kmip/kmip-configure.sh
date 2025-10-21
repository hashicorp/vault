#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
    echo "$1" 1>&2
    exit 1
}

[[ -z "${MOUNT}" ]] && fail "MOUNT env variable has not been set"
[[ -z "${KMIP_LISTEN_ADDR}" ]] && fail "KMIP_LISTEN_ADDR env variable has not been set"
[[ -z "${KMIP_PORT}" ]] && fail "KMIP_PORT env variable has not been set"
[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${VAULT_INSTALL_DIR}" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "${VAULT_TOKEN}" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "${binpath}" || fail "unable to locate vault binary at ${binpath}"

export VAULT_FORMAT=json

# Configure KMIP settings - redirect output to stderr to keep stdout clean
"${binpath}" write "${MOUNT}/config" \
    listen_addrs="${KMIP_LISTEN_ADDR}":"${KMIP_PORT}" \
    server_hostnames="${VAULT_ADDR}" >&2

# Read the CA certificate directly to stdout (no intermediate file needed)
"${binpath}" read "${MOUNT}"/ca -format=json | jq -r '.data | .ca_pem'
