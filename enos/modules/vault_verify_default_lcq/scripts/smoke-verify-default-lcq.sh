#!/bin/bash
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
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault

getMaxLeases() {
  $binpath read --format=json sys/quotas/lease-count/default | jq -r '.data.max_leases // empty'
}

waitForMaxLeases() {
  local max_leases
  if ! max_leases=$(getMaxLeases); then
    echo "failed getting /v1/sys/quotas/lease-count/default data" 1>&2
    return 1
  fi

  if [[ "$max_leases" == "$DEFAULT_LCQ" ]]; then
    echo "$max_leases"
    return 0
  else
    echo "Expected Default LCQ $DEFAULT_LCQ but got $max_leases"
    return 1
  fi
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if waitForMaxLeases; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for Default LCQ verification to complete. Data:\n\t$(getMaxLeases)\n"