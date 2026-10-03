#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

IFS="," read -r -a keys <<< "${UNSEAL_KEYS}"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

result=$($binpath operator unseal "${keys[0]}")
code=$?
if [ $code -eq 0 ]; then
  fail "expected unseal to fail but got exit code $code: $result"
fi
