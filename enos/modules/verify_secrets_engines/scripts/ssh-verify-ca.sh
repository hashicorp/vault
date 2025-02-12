#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$EXPECTED_CA_KEY_TYPE" ]] && fail "EXPECTED_CA_KEY_TYPE env variable has not been set"
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

# Extract actual key_type
actual_ca_key_type=$(echo "$ca_output" | jq -r '.data.key_type')

# Verify the key_type
[[ "$actual_ca_key_type" != "$EXPECTED_CA_KEY_TYPE" ]] && fail "CA key_type mismatch: expected $EXPECTED_CA_KEY_TYPE, got $actual_ca_key_type"

echo "SSH CA verification successful."
