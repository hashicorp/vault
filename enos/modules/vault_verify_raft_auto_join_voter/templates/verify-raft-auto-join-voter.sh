#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

binpath=${vault_install_dir}/vault

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
  voter_status=$($binpath operator raft list-peers -format json | jq -Mr --argjson expected "true" '.data.config.servers[] | select(.address=="${vault_cluster_addr}") | .voter == $expected')

  if [[ "$voter_status" != 'true' ]]; then
    fail "expected ${vault_cluster_addr} to be raft voter, got raft status for node: $($binpath operator raft list-peers -format json | jq '.data.config.servers[] | select(.address==${vault_cluster_addr})')"
  fi
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='${vault_token}'

# Retry a few times because it can take some time for things to settle after
# all the nodes are unsealed
retry 5 check_voter_status
