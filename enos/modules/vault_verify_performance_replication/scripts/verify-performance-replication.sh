#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  echo "$1" 2>&1
  return 1
}

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((10 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

check_pr_status() {
  cluster_state=$($binpath read -format=json sys/replication/performance/status | jq -r '.data.state')

  if [[ "${REPLICATION_MODE}" == "primary" ]]; then
    connection_status=$($binpath read -format=json sys/replication/performance/status | jq -r '.data.secondaries[0].connection_status')
  else
    connection_status=$($binpath read -format=json sys/replication/performance/status | jq -r '.data.primaries[0].connection_status')
  fi

  if [[ "$connection_status" == 'disconnected' ]]; then
    fail "expected connection status to be connected"
  fi

  if [[ "$cluster_state" == 'idle' ]]; then
    fail "expected cluster state to be not idle"
  fi
}

# Retry a few times because it can take some time for replication to sync
retry 5 check_pr_status
echo $($binpath read -format=json sys/replication/performance/status)
