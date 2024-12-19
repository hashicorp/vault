#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_LEADER_ADDR" ]] && fail "VAULT_LEADER_ADDR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

result=$($binpath operator raft join "$VAULT_LEADER_ADDR")
output=$?
if [ $output -ne 2 ]; then
  fail "Joining did not return code 2, instead $output: $result"  
fi
