#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$SET_NAME" ]] && fail "SET_NAME env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$MAX_TTL" ]] && fail "MAX_TTL env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath="${VAULT_INSTALL_DIR}/vault"
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

library_path="${MOUNT}/library/${SET_NAME}"

echo "Test Case #20: TTL Configuration for set=${SET_NAME}"

export VAULT_FORMAT=json

echo "Configuring TTL: ${TTL}, Max TTL: ${MAX_TTL}"

if output=$(
  "$binpath" write -format=json "$library_path" - << EOF 2>&1
{
  "ttl": "${TTL}",
  "max_ttl": "${MAX_TTL}"
}
EOF
); then
  printf "%s\n" "$output"
  echo "TTL configuration updated successfully"

  # Read back configuration to verify
  echo "Verifying TTL configuration:"
  read_output=$("$binpath" read -format=json "$library_path" 2>&1)

  configured_ttl=$(jq -r '.data.ttl' <<< "$read_output")
  configured_max_ttl=$(jq -r '.data.max_ttl' <<< "$read_output")

  echo "Configured TTL: ${configured_ttl}"
  echo "Configured Max TTL: ${configured_max_ttl}"
else
  fail "failed to configure TTL: set=${SET_NAME}, exit_code=$?, output: ${output}"
fi
