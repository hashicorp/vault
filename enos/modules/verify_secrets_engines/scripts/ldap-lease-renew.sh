#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LEASE_ID" ]] && fail "LEASE_ID env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo "Vault LEASE RENEW request for lease_id: $LEASE_ID"

if output=$("$binpath" lease renew "$LEASE_ID" 2>&1); then
  printf "%s\n" "$output"
else
  fail "failed to renew lease: lease_id=${LEASE_ID}, exit_code=$?, output: ${output}"
fi
