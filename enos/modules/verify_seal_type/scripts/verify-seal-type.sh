#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$EXPECTED_SEAL_TYPE" ]] && fail "EXPECTED_SEAL_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

count=0
retries=2
while :; do
  if seal_status=$($binpath read sys/seal-status -format=json); then
    if jq -Mer --arg expected "$EXPECTED_SEAL_TYPE" '.data.type == $expected' <<< "$seal_status" &> /dev/null; then
      exit 0
    fi
  fi

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else
    printf "Seal Status: %s\n" "$seal_status"
    got=$(jq -Mer '.data.type' <<< "$seal_status")
    fail "Expected seal type to be $EXPECTED_SEAL_TYPE, got: $got"
  fi
done
