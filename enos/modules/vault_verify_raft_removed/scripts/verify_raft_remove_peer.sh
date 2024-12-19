#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

getSysHealth() {
  curl -XGET --header "X-Vault-Token: $VAULT_TOKEN" \
    "$VAULT_ADDR/v1/sys/health" | jq '.removed_from_cluster'
}

getStatus() {
  $binpath status --format=json | jq '.removed_from_cluster'
}

expectRemoved() {
  local status
  status=$(getStatus)
  if [[ "$status" != "true" ]]; then
    echo "unexpected status $status"
    return 1
  fi
  
  local health
  health=$(getSysHealth)
  if [[ "$health" != "true" ]]; then
    echo "unexpected health $health"
  fi 
  
  return 0 
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if expectRemoved; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for raft removed status"
