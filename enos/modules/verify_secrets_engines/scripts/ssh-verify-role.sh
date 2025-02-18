#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$KEY_TYPE" ]] && fail "KEY_TYPE env variable has not been set"
[[ -z "$DEFAULT_USER" ]] && fail "DEFAULT_USER env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" read "ssh/roles/$ROLE_NAME" 2>&1); then
  fail "failed to read ssh/roles/$ROLE_NAME: $output"
fi

# Extract actual key type
key_type=$(echo "$output" | jq -r '.data.key_type')
default_user=$(echo "$output" | jq -r '.data.default_user')

# Verify
if [[ "$key_type" != "$KEY_TYPE" ]]; then
  fail "Key type mismatch: expected $KEY_TYPE, got $key_type"
fi

if [[ "$default_user" != "$DEFAULT_USER" ]]; then
  fail "Default user mismatch: expected $DEFAULT_USER, got $default_user"
fi
