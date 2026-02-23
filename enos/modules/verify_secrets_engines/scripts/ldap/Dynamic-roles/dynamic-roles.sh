#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

ROLE_NAME="${ROLE_NAME:-dynamic-role}"
# Test Case: create a dynamic role
test_create_dynamic_role() {
  echo "Test Case: Create Dynamic Role ($ROLE_NAME)"

  if ! output=$("$binpath" write "${MOUNT}/role/${ROLE_NAME}" \
      username_template="v-temp-{{random 10}}" \
      default_ttl=3200s \
      max_ttl=7200s 2>&1); then
    fail "ERROR: Failed to create Dynamic Role: ${output}"
  fi

  echo "SUCCESS: Dynamic Role created."
}

# Test Case: Read a newly created dynamic role
test_read_dynamic_role() {
  echo "Test Case: Read newly created dynamic role"

  if ! read_output=$("$binpath" read "${MOUNT}/role/${ROLE_NAME}" 2>&1); then
    fail "ERROR: Failed to read role: ${read_output}"
  fi

  if ! grep -q "v-temp" <<< "$read_output"; then
    echo "Debug Output: $read_output"
    fail "ERROR: Role read succeeded but creation_ldif not found."
  else
    echo "Debug Output: $read_output"
  fi
  echo "SUCCESS: Role is readable and configuration matches."
}

# Test Case: Update existing dynamic role
test_update_dynamic_role() {
  echo "Test case: Update existing dynamic role TTL"

  if ! output=$("$binpath" write "${MOUNT}/role/${ROLE_NAME}" default_ttl=7200s 2>&1); then
    fail "ERROR: Failed to update role: ${output}"
  fi

  if ! updated_ttl=$("$binpath" read -field=default_ttl "${MOUNT}/role/${ROLE_NAME}" 2>&1); then
    fail "ERROR: Failed to read updated TTL: ${updated_ttl}"
  fi

  if echo "$updated_ttl" | grep -q "7200"; then
    echo "SUCCESS: Role updated successfully (TTL=7200)."
  else
    echo "Debug Output: $updated_ttl"
    fail "ERROR: Update mismatch. Expected 7200, got $updated_ttl"
  fi
}

test_create_dynamic_role
test_read_dynamic_role
test_update_dynamic_role

echo "All basic dynamic role tests passed successfully."
