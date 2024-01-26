#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


binpath=${VAULT_INSTALL_DIR}/vault

function fail() {
  echo "$1" 1>&2
  exit 1
}

count=0
retries=5
while :; do
  # Check the Vault seal status
  seal_status=$($binpath status -format json | jq '.sealed')

  if [[ "$seal_status" == "true" ]]; then
    exit 0
  fi

  wait=$((3 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
  else
      fail "Expected node to be sealed"
  fi
done
