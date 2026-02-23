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

# Advanced LDAP dynamic role tests: Rollback, Deletion, Negative scenarios
# Test Case: Rollback on Creation Failure
test_rollback_on_creation_failure() {
  echo "Test Case: Rollback on Creation Failure (Negative Scenario)"

  FAIL_ROLE_NAME="rollback-test-role"
  FIXED_USER="rb_user_verify"

  cat << EOF > create_for_rollback.ldif
dn: uid={{.Username}},ou=users,dc=$LDAP_USERNAME,dc=com
objectClass: inetOrgPerson
sn: {{.Username}}
cn: {{.Username}}
uid: {{.Username}}
userPassword: {{.Password}}

dn: uid=fail-{{.Username}},ou=users,dc=$LDAP_USERNAME,dc=com
objectClass: THIS_CLASS_DOESNOT_EXIST
sn: fail-user
cn: fail-user
EOF

  cat << EOF > rollback.ldif
dn: uid={{.Username}},ou=users,dc=$LDAP_USERNAME,dc=com
changetype: delete
EOF

  "$binpath" write ldap/role/$FAIL_ROLE_NAME \
      creation_ldif=@create_for_rollback.ldif \
      deletion_ldif=@deletion.ldif \
      rollback_ldif=@rollback.ldif \
      username_template="$FIXED_USER" \
      default_ttl=3600s > /dev/null

  echo "Triggering credential generation for $FAIL_ROLE_NAME (Expected to FAIL)..."

  if "$binpath" read ldap/creds/$FAIL_ROLE_NAME > /dev/null 2>&1; then
    fail "ERROR: API succeeded unexpectedly. The invalid LDIF should have caused a failure."
  else
    echo "SUCCESS: API returned error as expected (simulated partial failure)."
  fi

  echo "Verifying Rollback: Checking if '$FIXED_USER' was properly cleaned up..."

  SEARCH_RES=$(ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" \
      -b "ou=users,dc=${LDAP_USERNAME},dc=com" \
      -D "cn=admin,dc=${LDAP_USERNAME},dc=com" \
      -w "${LDAP_ADMIN_PW}" \
      "(uid=$FIXED_USER)" 2>&1)

  if grep -q "^dn:" <<< "$SEARCH_RES"; then
    fail "ERROR: user $FIXED_USER found in ldap."
  else
    echo "SUCCESS: User not found! Rollback successful"
  fi
}

test_rollback_on_creation_failure

echo "Rollback tests passed successfully."
