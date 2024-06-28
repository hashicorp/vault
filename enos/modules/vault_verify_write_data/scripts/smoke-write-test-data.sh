#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

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
    else
      return "$exit"
    fi
  done

  return 0
}

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$TEST_KEY" ]] && fail "TEST_KEY env variable has not been set"
[[ -z "$TEST_VALUE" ]] && fail "TEST_VALUE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

retry 5 "$binpath" kv put secret/test "$TEST_KEY=$TEST_VALUE"
