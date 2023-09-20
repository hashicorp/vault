#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

# shellcheck disable=SC2154
binpath=${vault_install_dir}/vault

fail() {
  echo "$1" 1>&2
  exit 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'

count=0
retries=4
while :; do
    health_status=$(curl http://127.0.0.1:8200/v1/sys/health |jq '.')
    unseal_status=$($binpath status -format json | jq -Mr --argjson expected "false" '.sealed == $expected')
    if [[ "$unseal_status" == 'true' ]]; then
      echo "$health_status"
      exit 0
    fi

    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      # shellcheck disable=SC2154
      fail "expected ${vault_cluster_addr} to be unsealed, got unseal status: $unseal_status"
    fi
done
