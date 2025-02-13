#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$CA_KEY_TYPE" ]] && fail "CA_KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# Read the SSH CA configuration from Vault
if ! ca_output=$("$binpath" read "ssh/config/ca" 2>&1); then
  fail "failed to read ssh/config/ca: $ca_output"
fi

# Extract the CA public key
ca_public_key=$(echo "$ca_output" | jq -r '.data.public_key')

# Extract the first word (key type) from the public key
key_type=$(echo "$ca_public_key" | awk '{print $1}')

# Verify that the key type matches the expected key type
if [[ "$key_type" != "$CA_KEY_TYPE" ]]; then
  fail "CA key type mismatch: expected $CA_KEY_TYPE, got $key_type"
fi
