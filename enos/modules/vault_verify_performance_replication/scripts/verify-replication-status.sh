#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# This script waits for the replication status to be established
# then verifies the performance replication between primary and
# secondary clusters

set -e

binpath=${VAULT_INSTALL_DIR}/vault

function fail() {
	echo "$1" 1>&2
	exit 1
}

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done
}

test -x "$binpath" || exit 1

check_pr_status() {
  pr_status=$($binpath read -format=json sys/replication/performance/status)
  cluster_state=$(echo $pr_status | jq -r '.data.state')
  connection_mode=$(echo $pr_status | jq -r '.data.mode')

  if [[ "$cluster_state" == 'idle' ]]; then
    fail "replication cluster state is $cluster_state"
  fi

  if [[ "$connection_mode" == "primary" ]]; then
    connection_status=$(echo $pr_status | jq -r '.data.secondaries[0].connection_status')
    if [[ "$connection_status" == 'disconnected' ]]; then
      fail "replication connection status of secondaries is $connection_status"
    fi
    secondary_cluster_addr=$(echo $pr_status | jq -r '.data.secondaries[0].cluster_address')
    if [[ "$secondary_cluster_addr" != "https://"${SECONDARY_LEADER_PRIV_IP}":8201" ]]; then
      fail "Expected secondary cluster address $SECONDARY_LEADER_PRIV_IP got  $secondary_cluster_addr "
    fi
  else
    connection_status=$(echo $pr_status | jq -r '.data.primaries[0].connection_status')
    if [[ "$connection_status" == 'disconnected' ]]; then
      fail "replication connection status of secondaries is $connection_status"
    fi
    primary_cluster_addr=$(echo $pr_status | jq -r '.data.primaries[0].cluster_address')
    if [[ "$primary_cluster_addr" != "https://"${PRIMARY_LEADER_PRIV_IP}":8201" ]]; then
      fail "Expected primary cluster address $PRIMARY_LEADER_PRIV_IP got  $primary_cluster_addr"
    fi
    known_primary_cluster_addrs=$(echo $pr_status | jq -r '.data.known_primary_cluster_addrs')
    # IFS="," read -a cluster_addr <<< ${known_primary_cluster_addrs}
    if ! $(echo $known_primary_cluster_addrs |grep -q $PRIMARY_LEADER_PRIV_IP); then
      fail "Primary leader address $PRIMARY_LEADER_PRIV_IP not found in Known primary cluster addresses $known_primary_cluster_addrs"
    fi
  fi
  echo $pr_status
}

# Retry a few times because it can take some time for replication to sync
retry 5 check_pr_status
