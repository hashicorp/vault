#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$SET_NAME" ]] && fail "SET_NAME env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath="${VAULT_INSTALL_DIR}/vault"
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

checkout_path="${MOUNT}/library/${SET_NAME}/check-out"

echo "Test Case #17: Password Retrieval on Check-out for set=${SET_NAME}"

export VAULT_FORMAT=json

if output=$("$binpath" write -format=json -f "$checkout_path" 2>&1); then
  printf "%s\n" "$output"

  # Validate checkout response contains password
  service_account=$(jq -r '.data.service_account_name' <<< "$output")
  password=$(jq -r '.data.password' <<< "$output")
  lease_id=$(jq -r '.lease_id' <<< "$output")

  if [[ -z "$service_account" ]]; then
    fail "checkout response missing service_account_name: out=${output}"
  fi

  if [[ -z "$password" ]]; then
    fail "checkout response missing password: out=${output}"
  fi

  if [[ -z "$lease_id" ]]; then
    fail "checkout response missing lease_id: out=${output}"
  fi

  echo "Password retrieval validated:"
  echo "  Service Account: ${service_account}"
  echo "  Password Length: ${#password} characters"
  echo "  Lease ID: ${lease_id}"

  # Check in the account so it's available for any subsequent tests
  checkin_path="${MOUNT}/library/${SET_NAME}/check-in"
  if checkin_output=$(
    "$binpath" write -format=json "$checkin_path" - << EOF 2>&1
{
  "service_account_names": ["${service_account}"]
}
EOF
  ); then
    echo "  Checked in: ${service_account}"
  else
    echo "Warning: Check-in after test failed (non-fatal): ${checkin_output}"
  fi
else
  printf "%s\n" "$output" >&2
  fail "checkout failed: set=${SET_NAME}, exit_code=$?, output: ${output}"
fi
