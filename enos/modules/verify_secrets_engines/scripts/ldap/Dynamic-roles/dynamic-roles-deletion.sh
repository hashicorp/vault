#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}
retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
      echo "retry $count"
    else
      return "$exit"
    fi
  done

  return 0
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
# Test Case: Delete a dynamic role
test_delete_nonexistent_role() {
  echo "Test Case: Delete a nonexistent dynamic role"

  local test_role="${ROLE_NAME}_deletion"

  if "$binpath" read "${MOUNT}/role/${test_role}" > /dev/null 2>&1; then
    fail "ERROR: Role '$test_role' unexpectedly exists."
  fi

  echo "Confirmed: Role '$test_role' does not exist."

  "$binpath" delete "${MOUNT}/role/${test_role}" > /dev/null 2>&1

  if "$binpath" read "${MOUNT}/role/${test_role}" > /dev/null 2>&1; then
    fail "FAIL: Role unexpectedly exists after deletion command!"
  fi

  echo "PASS: Role deleted successfully. Read returned error as expected."
}
# Test Case: Role Deleted with Active Leases
test_role_deletion_with_active_leases() {
  echo "TEST CASE: Role Deleted with Active Leases (Cleanup)"

  local lease_role="${ROLE_NAME}_lease_cleanup"

  if ! output=$("$binpath" write "${MOUNT}/role/${lease_role}" \
      creation_ldif=@creation.ldif \
      deletion_ldif=@deletion.ldif \
      rollback_ldif=@rollback.ldif \
      username_template="v-lease-{{random 5}}" \
      default_ttl=3600s 2>&1); then
    fail "ERROR: Failed to create role: ${output}"
  fi

  echo "Role '$lease_role' created."

  if ! creds_output=$("$binpath" read -format=json "${MOUNT}/creds/${lease_role}" 2>&1); then
    fail "ERROR: Failed to generate credentials: ${creds_output}"
  fi

  lease_user=$(echo "$creds_output" | jq -r .data.username)

  if [[ -z "$lease_user" || "$lease_user" == "null" ]]; then
    fail "ERROR: Could not generate a user: ${creds_output}"
  fi

  echo "Generated dynamic user: '$lease_user'"

  "$binpath" delete "${MOUNT}/role/${lease_role}" > /dev/null 2>&1
  echo "Deleted role '$lease_role'. Waiting for cleanup..."

  "$binpath" lease revoke -prefix "${MOUNT}/creds/${lease_role}" > /dev/null 2>&1

  # Define the check function
  wait_for_user_deletion() {
    check_user=$(ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" \
        -b "dc=${LDAP_USERNAME},dc=com" \
        -D "cn=admin,dc=${LDAP_USERNAME},dc=com" \
        -w "${LDAP_ADMIN_PW}" \
        "(uid=$lease_user)" dn 2>&1)

    ! grep -q "^dn:" <<< "$check_user"
  }

  # Use retry to poll (6 attempts = ~63 seconds max)
  if retry 6 wait_for_user_deletion; then
    echo "SUCCESS: User '$lease_user' deleted!"
  else
    fail "ERROR: User still exists after cleanup timeout"
  fi
}

# Test Case: Invalid LDIF Template
test_invalid_ldif_template() {
  echo "--- TEST CASE: Invalid LDIF Template (Negative Input) ---"

   local bad_role="${ROLE_NAME}_invalid_ldif"

  echo "Attempting to write invalid LDIF..."

  set +e
  output=$("$binpath" write "${MOUNT}/role/${bad_role}" \
      creation_ldif="dn: this-is-total-garbage" \
      deletion_ldif=@delete.ldif 2>&1)
  ret_code=$?
  set -e

  if [ $ret_code -ne 0 ]; then
    if grep -iq "invalid\|failed to parse" <<< "$output"; then
      echo "SUCCESS: System correctly rejected invalid LDIF."
      echo "Error Message: ${output}"
    else
      echo "SUCCESS: System rejected invalid LDIF (generic error)."
      echo "Error Message: ${output}"
    fi
  else
    fail "FAIL: System accepted invalid data: ${output}"
  fi
}
