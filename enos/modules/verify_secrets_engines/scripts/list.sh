#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$REQPATH" ]] && fail "REQPATH env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo "Vault LIST request to path: $REQPATH"

set +e
output=$("$binpath" list "$REQPATH" 2>&1)
exit_code=$?
set -e

# Always print output
if [ "$exit_code" -eq 0 ]; then
  printf "%s\n" "$output"
else
  printf "%s\n" "$output" >&2
  # Exit code 2 typically means "not found" or "not listable" in vault
  if [ "$exit_code" -eq 2 ]; then
    echo "Note: Path is not listable or does not exist (exit code 2, this may be expected)" >&2
    exit 0
  fi
  # For other errors, fail the test
  fail "failed to list path: $REQPATH exit_code=${exit_code}"
fi
