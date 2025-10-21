#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

getStatus() {
  $binpath status -format json
}

isUnsealed() {
  local status
  if ! status=$(getStatus); then
    echo "failed to get vault status" 1>&2
    return 1
  fi

  if status=$(jq -Mre --argjson expected "false" '.sealed == $expected' <<< "$status"); then
    echo "vault is unsealed: $status"
    return 0
  fi

  echo "vault is sealed" 1>&2
  return 1
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  echo "waiting for vault to be unsealed..."

  if isUnsealed; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

if [ -n "$HOST_IPV6" ]; then
  fail "timed out waiting for Vault cluster on ${HOST_IPV6} to be unsealed"
fi
if [ -n "$HOST_IPV4" ]; then
  fail "timed out waiting for Vault cluster on ${HOST_IPV4} to be unsealed"
fi
fail "timed out waiting for Vault cluster to be unsealed"
