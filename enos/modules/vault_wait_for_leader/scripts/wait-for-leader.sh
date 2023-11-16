#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_INSTANCE_PRIVATE_IPS" ]] && fail "VAULT_INSTANCE_PRIVATE_IPS env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

findLeaderInPrivateIPs() {
  # Find the leader private IP address
  local leader_private_ip
  if ! leader_private_ip=$($binpath read sys/leader -format=json | jq -er '.data.leader_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")') ; then
    # Some older versions of vault don't support reading sys/leader. Fallback to the cli status.
    if ! leader_private_ip=$($binpath status -format json | jq -er '.leader_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")'); then
      return 1
    fi
  fi

  if isIn=$(jq -er --arg ip "$leader_private_ip" 'map(select(. == $ip)) | length == 1' <<< "$VAULT_INSTANCE_PRIVATE_IPS"); then
    if [[ "$isIn" == "true" ]]; then
      echo "$leader_private_ip"
      return 0
    fi
  fi

  return 1
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if findLeaderInPrivateIPs; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for one of $VAULT_INSTANCE_PRIVATE_IPS to be leader."
