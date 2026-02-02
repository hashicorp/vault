#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

# Skip if LEASE_ID is empty (lease not available from checkout)
if [[ -z "$LEASE_ID" ]] || [[ "$LEASE_ID" == "" ]]; then
  echo "Warning: LEASE_ID not set, skipping lease renew test"
  exit 0
fi

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo "Vault LEASE RENEW request for lease_id: $LEASE_ID"

set +e
output=$("$binpath" lease renew "$LEASE_ID" 2>&1)
exit_code=$?
set -e

# Always print output
if [ "$exit_code" -eq 0 ]; then
  printf "%s\n" "$output"
else
  printf "%s\n" "$output" >&2
  fail "failed to renew lease: lease_id=${LEASE_ID} exit_code=${exit_code}"
fi
