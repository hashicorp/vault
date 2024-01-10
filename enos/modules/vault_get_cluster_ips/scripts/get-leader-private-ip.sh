#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "Unable to locate vault binary at $binpath"

count=0
retries=5
while :; do
  # Find the leader private IP address
  if ip=$($binpath read sys/leader -format=json | jq -r '.data.leader_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")'); then
    if [[ -n "$ip" ]]; then
      echo "$ip"
      exit 0
    fi
  fi

  # Some older versions of vault don't support reading sys/leader. Try falling back to the cli status.
  if ip=$($binpath status -format json | jq -r '.leader_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")'); then
    if [[ -n "$ip" ]]; then
      echo "$ip"
      exit 0
    fi
  fi

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else
    fail "Timed out trying to obtain the cluster leader"
  fi
done
