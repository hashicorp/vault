#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$PAYLOAD" ]] && fail "PAYLOAD env variable has not been set"
[[ -z "$ASSERT_ACTIVE" ]] && fail "ASSERT_ACTIVE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" write identity/oidc/introspect - <<< "$PAYLOAD" 2>&1); then
  # Attempt to write our error on stdout as JSON as our consumers of the script expect it to be JSON
  printf '{"data":{"error":"%s"}}' "$output"
  # Fail on stderr with a human readable message
  fail "failed to write payload to identity/oidc/introspect: payload=$PAYLOAD output=$output"
fi

printf "%s\n" "$output" # Write our response output JSON to stdout
if ! jq -Me --argjson ACTIVE "$ASSERT_ACTIVE" '.data.active == $ACTIVE' <<< "$output" &> /dev/null; then
  # Write a failure message on STDERR
  fail "token active state is invalid, expected .data.active='$ASSERT_ACTIVE'"
fi
