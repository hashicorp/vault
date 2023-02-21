#!/usr/bin/env bash

set -e

binpath=${vault_install_dir}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'

health_status=$(curl http://127.0.0.1:8200/v1/sys/health |jq '.')
unseal_status=$($binpath status -format json | jq -Mr --argjson expected "false" '.sealed == $expected')
if [[ "$unseal_status" != 'true' ]]; then
  fail "expected ${vault_cluster_addr} to be unsealed, got unseal status: $unseal_status"
fi

echo $health_status
