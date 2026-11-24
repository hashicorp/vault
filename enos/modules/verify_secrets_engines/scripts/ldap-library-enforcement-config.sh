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
[[ -z "$SERVICE_ACCOUNT_NAMES" ]] && fail "SERVICE_ACCOUNT_NAMES env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$MAX_TTL" ]] && fail "MAX_TTL env variable has not been set"
[[ -z "$DISABLE_CHECK_IN_ENFORCEMENT" ]] && fail "DISABLE_CHECK_IN_ENFORCEMENT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath="${VAULT_INSTALL_DIR}/vault"
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

library_path="${MOUNT}/library/${SET_NAME}"

echo "Test Case #19: Optional Check-In Enforcement for set=${SET_NAME}"

export VAULT_FORMAT=json

echo "Configuring library with disable_check_in_enforcement=${DISABLE_CHECK_IN_ENFORCEMENT}"

# Convert comma-separated SERVICE_ACCOUNT_NAMES to JSON array
service_account_names_json=$(jq -RMec 'split(",")' <<< "$SERVICE_ACCOUNT_NAMES")

if output=$(
  "$binpath" write -format=json "$library_path" - << EOF 2>&1
{
  "service_account_names": ${service_account_names_json},
  "ttl": "${TTL}",
  "max_ttl": "${MAX_TTL}",
  "disable_check_in_enforcement": ${DISABLE_CHECK_IN_ENFORCEMENT}
}
EOF
); then
  printf "%s\n" "$output"
  echo " Library configured successfully with check-in enforcement setting"

  # Read back configuration to verify
  echo "Verifying configuration:"
  read_output=$("$binpath" read -format=json "$library_path" 2>&1)

  if jq -Merc --arg expected "$DISABLE_CHECK_IN_ENFORCEMENT" '.data.disable_check_in_enforcement as $got | ($expected | ascii_downcase) == ($got | tostring | ascii_downcase)' <<< "$read_output"; then
    echo "✓ disable_check_in_enforcement verified: ${DISABLE_CHECK_IN_ENFORCEMENT}"
  else
    fail "disable_check_in_enforcement mismatch: expected=${DISABLE_CHECK_IN_ENFORCEMENT}"
  fi
else
  fail "failed to configure library enforcement: set=${SET_NAME}, exit_code=$?, output: ${output}"
fi
