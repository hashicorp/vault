#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

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

check_voter_status() {
  voter_status=$($binpath operator raft list-peers -format json | jq -Mr --argjson expected "true" --arg ADDR "$VAULT_CLUSTER_ADDR" '.data.config.servers[] | select(.address==$ADDR) | .voter == $expected')

  if [[ "$voter_status" != 'true' ]]; then
    fail "expected $VAULT_CLUSTER_ADDR to be raft voter, got raft status for node: $($binpath operator raft list-peers -format json | jq -Mr --arg ADDR "$VAULT_CLUSTER_ADDR" '.data.config.servers[] | select(.address==$ADDR)')"
  fi
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

# Retry a few times because it can take some time for things to settle after
# all the nodes are unsealed
retry 10 check_voter_status
