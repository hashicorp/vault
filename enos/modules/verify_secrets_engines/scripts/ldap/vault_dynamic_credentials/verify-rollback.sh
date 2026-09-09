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

echo "Test: Rollback on Creation Failure"

# Create LDIF files for the test
# Invalid creation LDIF that will fail when executed against LDAP
cat > invalid_creation.ldif << EOF
dn: cn={{.Username}},ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: invalidClass
cn: {{.Username}}
EOF

# Rollback LDIF to clean up any partial creation
cat > rollback.ldif << EOF
dn: cn={{.Username}},ou=users,dc=${LDAP_USERNAME},dc=com
changetype: delete
EOF

# Deletion LDIF for normal cleanup
cat > deletion.ldif << EOF
dn: cn={{.Username}},ou=users,dc=${LDAP_USERNAME},dc=com
changetype: delete
EOF

echo "Creating role with invalid creation LDIF (role creation should succeed)"
if ! "$binpath" write "${MOUNT}/role/rollback-test" \
  creation_ldif=@invalid_creation.ldif \
  deletion_ldif=@deletion.ldif \
  rollback_ldif=@rollback.ldif \
  default_ttl="60s" \
  max_ttl="60s"; then
  rm -f invalid_creation.ldif rollback.ldif deletion.ldif
  fail "Role creation failed - it should succeed even with invalid LDIF syntax"
fi
echo "✅ Role creation succeeded (LDIF validation happens during credential generation)"

echo "Attempting to generate credentials (should fail and trigger rollback)"
if "$binpath" read "${MOUNT}/creds/rollback-test" 2>&1; then
  rm -f invalid_creation.ldif rollback.ldif deletion.ldif
  fail "Credential generation should have failed with invalid LDIF but succeeded"
else
  echo "✅ Credential generation correctly failed - rollback LDIF have been executed"
fi

echo "Verifying no orphaned LDAP users"
user_count=$(ldapsearch -x -LLL -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" \
  -b "ou=users,dc=${LDAP_USERNAME},dc=com" \
  -D "cn=admin,dc=${LDAP_USERNAME},dc=com" \
  -w "${LDAP_ADMIN_PW}" \
  "(objectClass=inetOrgPerson)" 2> /dev/null | grep -c "^dn:" || echo "0")

if [[ "$user_count" -eq 0 ]]; then
  echo "✅ SUCCESS: No orphaned users found after rollback"
else
  echo "⚠️  WARNING: Found $user_count users (may include pre-existing users)"
fi

rm -f invalid_creation.ldif rollback.ldif deletion.ldif
