#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

# Validate required environment variables
[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$STRONG_POLICY" ]] && fail "STRONG_POLICY env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
WEAK_POLICY="weak-policy"

# Test: Password policy enforcement with enos user
echo "Test: Creating weak password policy with min length 4"

cat > "${WEAK_POLICY}.hcl" << EOF
length = 4

rule "charset" {
  charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  min-chars = 1
}
EOF

if policy_output=$("$binpath" write sys/policies/password/"${WEAK_POLICY}" policy=@"${WEAK_POLICY}.hcl" 2>&1); then
  echo "Weak password policy created"
else
  fail "ERROR: Failed to create weak password policy $policy_output"
fi

# Configure LDAP with enos user and weak policy
echo "Vault: Configuring LDAP secrets engine with enos user"

if config_output=$("$binpath" write "${MOUNT}/config" \
  binddn="uid=${LDAP_USERNAME},ou=users,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  password_policy="$WEAK_POLICY" 2>&1); then
  echo "LDAP config updated with enos user and weak policy"
else
  fail "ERROR: Failed to update LDAP config with weak policy $config_output"
fi

# Attempt to rotate root credentials for enos user
echo "Test: Attempting to rotate root credentials for enos user"
if rotation_output=$("$binpath" write -f "${MOUNT}/rotate-root" 2>&1); then
  fail "ERROR: Rotate root should have failed with password policy constraint violation $rotation_output"
else
  echo "SUCCESS: Rotate root failed as expected due to password policy constraint"
fi

# Verify enos user credentials are still valid after rotation: should be valid, as rotation failed.
echo "Test: Verifying enos user credentials after rotation"
if ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "uid=${LDAP_USERNAME},ou=users,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" > /dev/null 2>&1; then
  echo "Old password valid after rotation failed"
else
  fail "Old password should have been valid after rotation failure"
fi

# Revert back to earlier password policy and admin user
echo "Test: Reverting LDAP config back to strong-policy and admin user"

"$binpath" write "${MOUNT}/config" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  password_policy="$STRONG_POLICY"

echo "SUCCESS: LDAP config reverted to strong-policy with admin user"
