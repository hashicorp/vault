#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

log() {
  echo "[DEBUG] $1" >&2
}

[[ -z "$VERIFY_SSH_SECRETS" ]] && fail "VERIFY_SSH_SECRETS env variable has not been set"
[[ -z "$CA_KEY_TYPE" ]] && fail "CA_KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

if [[ "$VERIFY_SSH_SECRETS" == "false" ]]; then
  log "VERIFY_SSH_SECRETS is false; exiting script"
  exit 0
fi

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# Read the SSH CA configuration from Vault
if ! ca_output=$("$binpath" read "ssh/config/ca" 2>&1); then
  fail "failed to read ssh/config/ca: $ca_output"
fi

# Extract the CA public key
ca_public_key=$(jq -r '.data.public_key' <<< "$ca_output")

# Extract the first word (key type) from the public key
key_type=$(awk '{print $1}' <<< "$ca_public_key")

# Verify that the key type matches the expected key type
if [[ "$key_type" != "$CA_KEY_TYPE" ]]; then
  fail "CA key type mismatch: expected $CA_KEY_TYPE, got $key_type"
fi
