#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"


while :; do
  default_lcq_resp=$(curl --request GET --header "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR"/v1/sys/quotas/lease-count/default)
  max_leases=$(jq '.data.max_leases // empty' <<< "$default_lcq_resp")
  if [[ "$max_leases" == "${DEFAULT_LCQ}" ]]; then
    exit 0
  else
    echo "Expected Default LCQ $DEFAULT_LCQ but got $max_leases"
    exit 1
  fi

done
