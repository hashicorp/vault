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

# Test Case: Validate dynamic role Schema
test_validate_schema() {
  echo "Test: Validate dynamic role schema"

  if ! output=$("$binpath" write "${MOUNT}/role/${ROLE_NAME}_schema" \
      creation_ldif=@creation.ldif \
      deletion_ldif=@deletion.ldif \
      rollback_ldif=@rollback.ldif \
      username_template="v-schema-{{random 5}}" 2>&1); then
    fail "ERROR: Schema validation failed: ${output}"
  fi

  echo "SUCCESS: Schema validation successful."
}

# Test Case: Listing Dynamic Roles
test_list_dynamic_roles() {
  echo "Test: List Dynamic Roles"

  if ! list_output=$("$binpath" list "${MOUNT}/role" 2>&1); then
    fail "ERROR: Failed to list roles: ${list_output}"
  fi

  if ! grep -q "$ROLE_NAME" <<< "$list_output"; then
    fail "ERROR: Role not found in list: ${list_output}"
  fi

  echo "SUCCESS: Role '$ROLE_NAME' found in list."
  echo "$list_output"
}

test_validate_schema
test_list_dynamic_roles

echo "All validation and listing tests passed successfully."
