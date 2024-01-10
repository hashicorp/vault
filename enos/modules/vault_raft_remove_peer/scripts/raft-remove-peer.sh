#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

binpath=${VAULT_INSTALL_DIR}/vault
node_addr=${REMOVE_VAULT_CLUSTER_ADDR}

fail() {
  echo "$1" 2>&1
  return 1
}

[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

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
      echo "retry $count"
    else
      return "$exit"
    fi
  done

  return 0
}

remove_peer() {
  if ! node_id=$("$binpath" operator raft list-peers -format json | jq -Mr --argjson expected "false" '.data.config.servers[] | select(.address=='\""$node_addr"\"') | select(.voter==$expected) | .node_id'); then
    fail "failed to get node id of a non-voter node"
  fi

  $binpath operator raft remove-peer "$node_id"
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# Retry a few times because it can take some time for things to settle after autopilot upgrade
retry 5 remove_peer
