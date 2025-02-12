#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

# Check required environment variables
[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$KEY_TYPE" ]] && fail "KEY_TYPE env variable has not been set"
[[ -z "$DEFAULT_USER" ]] && fail "DEFAULTUSER env variable has not been set"
[[ -z "$PORT" ]] && fail "PORT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" read "ssh/roles/$ROLE_NAME" 2>&1); then
  fail "failed to read ssh/roles/$ROLE_NAME: $output"
fi

# Extract actual values
actual_key_type=$(echo "$output" | jq -r '.data.key_type')
actual_user=$(echo "$output" | jq -r '.data.default_user')
actual_port=$(echo "$output" | jq -r '.data.port')

# Verify the values
[[ "$actual_key_type" != "$KEY_TYPE" ]] && fail "key_type mismatch: expected $KEY_TYPE, got $actual_key_type"
[[ "$actual_user" != "$DEFAULT_USER" ]] && fail "default_user mismatch: expected $DEFAULT_USER, got $actual_user"
[[ "$actual_port" != "$PORT" ]] && fail "port mismatch: expected $PORT, got $actual_port"

echo "SSH role verification successful."
