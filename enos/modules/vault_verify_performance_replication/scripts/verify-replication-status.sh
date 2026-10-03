#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# This script waits for the replication status to be established
# then verifies the performance replication between primary and
# secondary clusters

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$IP_VERSION" ]] && fail "IP_VERSION env variable has not been set"
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
  cluster_state=$(jq -r '.data.state' <<< "$pr_status")
  connection_mode=$(jq -r '.data.mode' <<< "$pr_status")

  if [[ "$cluster_state" == 'idle' ]]; then
    echo "replication cluster state is idle" 1>&2
    return 1
  fi

  if [[ "$connection_mode" == "primary" ]]; then
    if [[ "$cluster_state" != "running" ]]; then
      echo "primary cluster is not running" 1>&2
      return 1
    fi
    connection_status=$(jq -r '.data.secondaries[0].connection_status' <<< "$pr_status")
    if [[ "$connection_status" == 'disconnected' ]]; then
      echo ".data.secondaries[0].connection_status from primary node is 'disconnected'" 1>&2
      return 1
    fi
  else
    if [[ "$cluster_state" != "stream-wals" ]]; then
      echo "secondary cluster is not streaming wals" 1>&2
      return 1
    fi
    connection_status=$(jq -r '.data.primaries[0].connection_status' <<< "$pr_status")
    if [[ "$connection_status" == 'disconnected' ]]; then
      echo ".data.primaries[0].connection_status from secondary node is 'disconnected'" 1>&2
      return 1
    fi
    if [ "$IP_VERSION" == 4 ]; then
      primary_cluster_addr=$(jq -r '.data.primaries[0].cluster_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")' <<< "$pr_status")
    else
      primary_cluster_addr=$(jq -r '.data.primaries[0].cluster_address | scan("\\[(.+)\\]") | .[0]' <<< "$pr_status")
    fi
    known_primary_cluster_addrs=$(jq -r '.data.known_primary_cluster_addrs' <<< "$pr_status")
    if ! echo "$known_primary_cluster_addrs" | grep -q "$primary_cluster_addr"; then
      echo "$primary_cluster_addr is not in .data.known_primary_cluster_addrs: $known_primary_cluster_addrs" 1>&2
      return 1
    fi
  fi

  echo "$pr_status"
  return 0
}

if [ "$IP_VERSION" != 4 ] && [ "$IP_VERSION" != 6 ]; then
  fail "unsupported IP_VERSION: $IP_VERSION"
fi

# Retry for a while because it can take some time for replication to sync
retry 10 check_pr_status
