#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  return 1
}

[[ -z "$REPLICATION_TYPE" ]] && fail "REPLICATION_TYPE env variable has not been set"
[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

replicationStatus() {
  $binpath read "sys/replication/${REPLICATION_TYPE}/status" | jq .data
}

isReady() {
  # Find the leader private IP address
  local status
  if ! status=$(replicationStatus); then
    return 1
  fi

  if ! jq -eMc '.state == "stream-wals"' &> /dev/null <<< "$status"; then
    echo "DR replication state is not yet running" 1>&2
    echo "DR replication is not yet running, got: $(jq '.state' <<< "$status")" 1>&2
    return 1
  fi

  if ! jq -eMc '.mode == "secondary"' &> /dev/null <<< "$status"; then
    echo "DR replication mode is not yet primary, got: $(jq '.mode' <<< "$status")" 1>&2
    return 1
  fi

  if ! jq -eMc '.corrupted_merkle_tree == false' &> /dev/null <<< "$status"; then
    echo "DR replication merkle is corrupted" 1>&2
    return 1
  fi

  echo "${REPLICATION_TYPE} primary is ready for followers to be unsealed!" 1>&2
  return 0
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if isReady; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for ${REPLICATION_TYPE} primary to ready: $(replicationStatus)"
