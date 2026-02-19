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

# Test Case: Template Validation - Valid Syntax
test_valid_template_syntax() {
  echo "Test Case: Template Validation (Valid Syntax)"

  if ! output=$("$binpath" write "${MOUNT}/role/${ROLE_NAME}_valid" \
      creation_ldif=@creation.ldif \
      deletion_ldif=@deletion.ldif \
      rollback_ldif=@rollback.ldif \
      username_template="v-temp-{{random 10}}" 2>&1); then
    fail "ERROR: Vault rejected valid template at role creation: ${output}"
  fi

  echo "SUCCESS: Vault accepted valid template (200 OK)"
}

# Type 2: Negative Test - Invalid Template
test_invalid_template_syntax() {
  echo "Test Case: Template Validation (Invalid Syntax)"

  set +e
  output=$("$binpath" write "${MOUNT}/role/${ROLE_NAME}_invalid" \
      creation_ldif=@creation.ldif \
      deletion_ldif=@deletion.ldif \
      rollback_ldif=@rollback.ldif \
      username_template="{{.Invalid_Syntax}}" 2>&1)
  set +e
  output=$("$binpath" read "${MOUNT}/creds/${ROLE_NAME}_invalid" 2>&1)
  ret_code=$?
  set -e

  if [ $ret_code -eq 0 ]; then
    fail "ERROR: Vault accepted invalid template and generated credentials! Output: ${output}"
  fi

  if echo "$output" | grep -q "can't evaluate field Invalid_Syntax"; then
    echo "SUCCESS: Vault correctly rejected invalid template at credential generation"
    echo "Error message: ${output}"
  else
    echo "Debug Output: ${output}"
    fail "ERROR: Expected template validation error but got different error"
  fi
}

test_valid_template_syntax
test_invalid_template_syntax

echo "All template validation tests passed successfully."
