#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# This script waits for the replication status to be established
# then verifies the dr replication between primary and
# secondary clusters

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$IP_VERSION" ]] && fail "IP_VERSION env variable has not been set"
[[ -z "$PRIMARY_LEADER_ADDR" ]] && fail "PRIMARY_LEADER_ADDR env variable has not been set"
[[ -z "$SECONDARY_LEADER_ADDR" ]] && fail "SECONDARY_LEADER_ADDR env variable has not been set"
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
      fail "$($binpath read -format=json sys/replication/dr/status)"
    fi
  done
}

check_dr_status() {
  dr_status=$($binpath read -format=json sys/replication/dr/status)
  cluster_state=$(jq -r '.data.state' <<< "$dr_status")
  connection_mode=$(jq -r '.data.mode' <<< "$dr_status")

  if [[ "$cluster_state" == 'idle' ]]; then
    echo "replication cluster state is idle" 1>&2
    return 1
  fi

  if [[ "$connection_mode" == "primary" ]]; then
    connection_status=$(jq -r '.data.secondaries[0].connection_status' <<< "$dr_status")
    if [[ "$connection_status" == 'disconnected' ]]; then
      echo ".data.secondaries[0].connection_status from primary node is 'disconnected'" 1>&2
      return 1
    fi
    # Confirm we are in a "running" state for the primary
    if [[ "$cluster_state" != "running" ]]; then
      echo "replication cluster primary state is not running" 1>&2
      return 1
    fi
  else
    connection_status=$(jq -r '.data.primaries[0].connection_status' <<< "$dr_status")
    if [[ "$connection_status" == 'disconnected' ]]; then
      echo ".data.primaries[0].connection_status from secondary node is 'disconnected'" 1>&2
      return 1
    fi
    # Confirm we are in a "stream-wals" state for the secondary
    if [[ "$cluster_state" != "stream-wals" ]]; then
      echo "replication cluster primary state is not stream-wals" 1>&2
      return 1
    fi
    known_primary_cluster_addrs=$(jq -r '.data.known_primary_cluster_addrs' <<< "$dr_status")
    if ! echo "$known_primary_cluster_addrs" | grep -q "$PRIMARY_LEADER_ADDR"; then
      echo "$PRIMARY_LEADER_ADDR is not in .data.known_primary_cluster_addrs: $known_primary_cluster_addrs" 1>&2
      return 1
    fi
  fi

  echo "$dr_status"
  return 0
}

if [ "$IP_VERSION" != 4 ] && [ "$IP_VERSION" != 6 ]; then
  fail "unsupported IP_VERSION: $IP_VERSION"
fi

# Retry for a while because it can take some time for replication to sync
retry 10 check_dr_status
