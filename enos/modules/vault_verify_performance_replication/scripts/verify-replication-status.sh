#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


# This script waits for the replication status to be established
# then verifies the performance replication between primary and
# secondary clusters

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$PRIMARY_LEADER_PRIV_IP" ]] && fail "PRIMARY_LEADER_PRIV_IP env variable has not been set"
[[ -z "$SECONDARY_LEADER_PRIV_IP" ]] && fail "SECONDARY_LEADER_PRIV_IP env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      fail "$($binpath read -format=json sys/replication/performance/status)"
    fi
  done
}

check_pr_status() {
  pr_status=$($binpath read -format=json sys/replication/performance/status)
  cluster_state=$(echo "$pr_status" | jq -r '.data.state')
  connection_mode=$(echo "$pr_status" | jq -r '.data.mode')

  if [[ "$cluster_state" == 'idle' ]]; then
    echo "replication cluster state is idle" 1>&2
    return 1
  fi

  if [[ "$connection_mode" == "primary" ]]; then
    connection_status=$(echo "$pr_status" | jq -r '.data.secondaries[0].connection_status')
    if [[ "$connection_status" == 'disconnected' ]]; then
      echo ".data.secondaries[0].connection_status from primary node is 'disconnected'" 1>&2
      return 1
    fi
    secondary_cluster_addr=$(echo "$pr_status" | jq -r '.data.secondaries[0].cluster_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")')
    if [[ "$secondary_cluster_addr" != "$SECONDARY_LEADER_PRIV_IP" ]]; then
      echo ".data.secondaries[0].cluster_address should have an IP address of $SECONDARY_LEADER_PRIV_IP, got: $secondary_cluster_addr" 1>&2
      return 1
    fi
  else
    connection_status=$(echo "$pr_status" | jq -r '.data.primaries[0].connection_status')
    if [[ "$connection_status" == 'disconnected' ]]; then
      echo ".data.primaries[0].connection_status from secondary node is 'disconnected'" 1>&2
      return 1
    fi
    primary_cluster_addr=$(echo "$pr_status" | jq -r '.data.primaries[0].cluster_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")')
    if [[ "$primary_cluster_addr" != "$PRIMARY_LEADER_PRIV_IP" ]]; then
      echo ".data.primaries[0].cluster_address should have an IP address of $PRIMARY_LEADER_PRIV_IP, got: $primary_cluster_addr" 1>&2
      return 1
    fi
    known_primary_cluster_addrs=$(echo "$pr_status" | jq -r '.data.known_primary_cluster_addrs')
    if ! echo "$known_primary_cluster_addrs" | grep -q "$PRIMARY_LEADER_PRIV_IP"; then
      echo "$PRIMARY_LEADER_PRIV_IP is not in .data.known_primary_cluster_addrs: $known_primary_cluster_addrs" 1>&2
      return 1
    fi
  fi

  echo "$pr_status"
  return 0
}

# Retry for a while because it can take some time for replication to sync
retry 10 check_pr_status
