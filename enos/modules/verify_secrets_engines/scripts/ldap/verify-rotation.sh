#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

# Function to perform root rotation
rotate_root() {
  "$binpath" write -f "${MOUNT}/rotate-root" 2>&1
}

# Function to get userPassword from LDAP
get_ldap_password() {
  local user_dn="$1"
  ldapsearch -x -LLL -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" \
    -b "${user_dn}" \
    -D "cn=admin,dc=${LDAP_USERNAME},dc=com" \
    -w "${LDAP_ADMIN_PW}" userPassword 2> /dev/null | grep "userPassword::" | awk '{print $2}'
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
ROLE_NAME="dynamic-role"

# Verifying LDAP Server Configs
LDAP_UID=$(ldapsearch -x -LLL -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=${LDAP_USERNAME})" 2> /dev/null)
[[ -z "$LDAP_UID" ]] && fail "Could not search ldap server for uid: ${LDAP_USERNAME}"

# Authenticate Using Vault LDAP login
VAULT_LDAP_LOGIN=$("$binpath" login -method="${MOUNT}" username="${LDAP_USERNAME}" password="${LDAP_ADMIN_PW}")

# Verifying Vault LDAP Login Token
VAULT_LDAP_TOKEN=$(echo "$VAULT_LDAP_LOGIN" | jq -r ".auth.client_token")
[[ -z "$VAULT_LDAP_TOKEN" ]] && fail "Vault LDAP could not log in correctly: ${VAULT_LDAP_TOKEN}"

# Verifying Vault LDAP Policies
VAULT_POLICY_COUNT=$(echo "$VAULT_LDAP_LOGIN" | jq -r ".auth.policies | length")
[[ -z "$VAULT_POLICY_COUNT" ]] && fail "Vault LDAP number of policies does not look correct: ${VAULT_POLICY_COUNT}"

echo "${VAULT_LDAP_LOGIN}"

# Test1: Attempting to rotate root with root token--should pass
test_root_rotation_permissions() {
  if rotate_root > /dev/null; then
    echo "SUCCESS: rotate-root succeeded"
  else
    fail "Error: rotate-root write failed even though token had permissions"
  fi
}

# Test2: Checking if last_bind_password_rotation field is present
test_rotation_field_presence() {
  if "$binpath" read "${MOUNT}/config" | jq -e '.data.last_bind_password_rotation' > /dev/null; then
    echo "Rotation success: last_bind_password_rotation field is present"
  else
    fail "Field is NOT present"
  fi
}

# Test3: Attempting to rotate root with LDAP token--should fail as the policy does not allow it
test_ldap_token_permissions() {
  if ! VAULT_TOKEN="$VAULT_LDAP_TOKEN" "$binpath" write -f "${MOUNT}/rotate-root" > /dev/null 2>&1; then
    echo "SUCCESS: Vault correctly denied root rotation for LDAP token as policy does not allow."
  else
    fail "ERROR: LDAP token does not have permission to rotate, still rotation succeeded"
  fi
}

# Test4: Rotation with Invalid Config
test_invalid_config_rotation() {
  echo "Test 4: Rotation with Invalid Config"

  # Get password before attempting rotation with invalid config
  PASSWORD_BEFORE=$(get_ldap_password "cn=admin,dc=${LDAP_USERNAME},dc=com")

  if [[ -z "$PASSWORD_BEFORE" ]]; then
    fail "ERROR: Could not retrieve password before rotation attempt"
  fi

  # Attempt to configure with invalid credentials
  "$binpath" write "${MOUNT}/config" \
    binddn="cn=invalid,dc=invalid,dc=com" \
    bindpass="wrongpassword" > /dev/null 2>&1 || true

  # Try to rotate with invalid config: should fail
  if ! rotate_root > /dev/null 2>&1; then
    echo "SUCCESS: Rotation correctly failed with invalid configuration"
  else
    fail "ERROR: Rotation should have failed with invalid configuration"
  fi

  # Restore valid config first
  "$binpath" write "${MOUNT}/config" \
    binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
    bindpass="${LDAP_ADMIN_PW}"

  # Get password after failed rotation attempt
  PASSWORD_AFTER=$(get_ldap_password "cn=admin,dc=${LDAP_USERNAME},dc=com")

  if [[ -z "$PASSWORD_AFTER" ]]; then
    fail "ERROR: Could not retrieve password after rotation attempt"
  fi

  # Verify password remains unchanged
  if [[ "$PASSWORD_BEFORE" == "$PASSWORD_AFTER" ]]; then
    echo "SUCCESS: User password remains unchanged after failed rotation with invalid config"
  else
    fail "ERROR: User password was modified despite rotation failure"
  fi

  # Verify credentials work for LDAP required operations by testing dynamic credential generation
  if "$binpath" read "${MOUNT}/creds/${ROLE_NAME}" > /dev/null 2>&1; then
    echo "SUCCESS: LDAP operations work correctly after restoring valid config"
  else
    fail "ERROR: LDAP operations failed after restoring valid config"
  fi
}

# Test5: Rotate root twice and check if password changed
test_password_change_verification() {
  echo "Performing first root rotation"
  rotate_root > /dev/null

  # Get password after first rotation, this will return the userPassword of the binddn.
  echo "Getting password after first rotation"
  FIRST_ROTATED_PASSWORD=$(get_ldap_password "cn=admin,dc=${LDAP_USERNAME},dc=com")

  if [[ -z "$FIRST_ROTATED_PASSWORD" ]]; then
    fail "ERROR: Could not retrieve password after first rotation from LDAP"
  fi
  echo "First rotated password retrieved"

  # Second rotation
  echo "Performing second root rotation"
  rotate_root > /dev/null

  # Get password after second rotation
  echo "Getting password after second rotation"
  SECOND_ROTATED_PASSWORD=$(get_ldap_password "cn=admin,dc=${LDAP_USERNAME},dc=com")

  if [[ -z "$SECOND_ROTATED_PASSWORD" ]]; then
    fail "ERROR: Could not retrieve password after second rotation from LDAP"
  fi

  # Compare passwords to ensure Userpassword is different after rotation
  if [[ "$FIRST_ROTATED_PASSWORD" == "$SECOND_ROTATED_PASSWORD" ]]; then
    fail "ERROR: Second rotation did not change the password! First and second rotated passwords are the same."
  fi

}
# Running all tests
test_root_rotation_permissions
test_rotation_field_presence
test_ldap_token_permissions
test_invalid_config_rotation
test_password_change_verification

echo "All rotation tests passed successfully!"
