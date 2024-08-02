#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e


fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

count=0
retries=5
while :; do
    health_status=$(curl -s "${VAULT_ADDR}/v1/sys/health" | jq '.')
    if unseal_status=$($binpath status -format json | jq -Mre --argjson expected "false" '.sealed == $expected'); then
      echo "$health_status"
      exit 0
    fi

    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      if [ -n "$HOST_IPV6" ]; then
        fail "expected ${HOST_IPV6} to be unsealed, got unseal status: $unseal_status"
      else
        if [ -n "$HOST_IPV4" ]; then
          fail "expected ${HOST_IPV4} to be unsealed, got unseal status: $unseal_status"
        else
          fail "expected ${VAULT_ADDR} to be unsealed, got unseal status: $unseal_status"
        fi
      fi
    fi
done
