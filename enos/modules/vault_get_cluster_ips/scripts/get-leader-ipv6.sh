#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$IP_VERSION" ]] && fail "IP_VERSION env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "Unable to locate vault binary at $binpath"

findLeaderIPV6() {
  # Find the leader private IP address
  if ip=$($binpath read sys/leader -format=json | jq -r '.data.leader_address | scan("\\[(.+)\\]") | .[0]'); then
    if [[ -n "$ip" ]]; then
      echo "$ip"
      return 0
    fi
  fi

  # Some older versions of vault don't support reading sys/leader. Try falling back to the cli status.
  if ip=$($binpath status -format json | jq -r '.leader_address | scan("\\[(.+)\\]") | .[0]'); then
    if [[ -n "$ip" ]]; then
      echo "$ip"
      return 0
    fi
  fi

  return 1
}

count=0
retries=5
while :; do
  # Find the leader private IP address
  case $IP_VERSION in
    4)
      exit 0
      ;;
    6)
      if ip=$(findLeaderIPV6); then
        echo "$ip"
        exit 0
      fi
      ;;
    *)
      fail "unknown IP_VERSION: $IP_VERSION"
      ;;
  esac

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else
    fail "Timed out trying to obtain the cluster leader"
  fi
done
