#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

# Function to generate LDIF files with DN
generate_ldif() {
  local filename="$1"
  local content="$2"
  local ou="${3:-users}"

  cat << EOF > "${filename}"
dn: cn={{.Username}},ou=${ou},dc=${LDAP_USERNAME},dc=com
${content}
EOF
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$DEFAULT_TTL" ]] && fail "DEFAULT_TTL env variable has not been set"
[[ -z "$MAX_TTL" ]] && fail "MAX_TTL env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo "=== Negative Test: LDAP Server Unreachable ==="

# Use unique mount path to avoid conflicts with parallel test runs
unreachable_mount="${UNREACHABLE_MOUNT:-ldap-unreachable-$(date +%s)}"
unreachable_role_name="${UNREACHABLE_ROLE_NAME:-test-unreachable}"
unreachable_server="${UNREACHABLE_SERVER:-unreachable-server.invalid}"
unreachable_port="${UNREACHABLE_PORT:-389}"

echo "Enabling separate LDAP mount for unreachable server test..."
if ! "$binpath" secrets enable -path="${unreachable_mount}" ldap 2>&1; then
  fail "Failed to enable LDAP mount at '${unreachable_mount}'"
fi

echo "Configuring LDAP mount with unreachable server..."
if ! "$binpath" write "${unreachable_mount}/config" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  url="ldap://${unreachable_server}:${unreachable_port}" 2>&1; then
  "$binpath" secrets disable "${unreachable_mount}" 2>&1 || echo "Warning: Failed to disable mount during cleanup"
  fail "Failed to configure LDAP mount with unreachable server"
fi

unreachable_creation="unreachable_creation.ldif"
unreachable_deletion="unreachable_deletion.ldif"
generate_ldif "${unreachable_creation}" "objectClass: person
objectClass: top
cn: {{.Username}}
sn: {{.Password | utf16le | base64}}
userPassword: {{.Password}}" "users"

generate_ldif "${unreachable_deletion}" "changetype: delete" "users"

echo "Creating role with valid LDIF (should succeed even with unreachable server)..."
if ! "$binpath" write "${unreachable_mount}/role/${unreachable_role_name}" \
  creation_ldif=@${unreachable_creation} \
  deletion_ldif=@${unreachable_deletion} \
  default_ttl="${DEFAULT_TTL}" \
  max_ttl="${MAX_TTL}" 2>&1; then
  rm -f "$unreachable_creation" "$unreachable_deletion"
  "$binpath" secrets disable "${unreachable_mount}" 2>&1 || echo "Warning: Failed to disable mount during cleanup"
  fail "Failed to create role - role creation should succeed even with unreachable server"
fi

echo "Attempting credential generation with unreachable LDAP server (should fail)..."
if "$binpath" read "${unreachable_mount}/creds/${unreachable_role_name}" 2>&1; then
  rm -f "$unreachable_creation" "$unreachable_deletion"
  "$binpath" secrets disable "${unreachable_mount}" 2>&1 || echo "Warning: Failed to disable mount during cleanup"
  fail "Vault should have failed when LDAP server is unreachable, but credentials were generated"
else
  echo "✅ SUCCESS: Vault correctly failed credential request with unreachable server"
fi

# Verify no lease was created by checking if the lease path exists
echo ""
echo "Verifying no phantom leases were created..."
if LEASE_LIST=$("$binpath" list "sys/leases/lookup/${unreachable_mount}/creds" 2>&1); then
  # If list succeeds, check if there are any leases
  if jq -e '.data.keys' <<< "$LEASE_LIST" > /dev/null 2>&1; then
    LEASE_COUNT=$(jq -r '.data.keys | length' <<< "$LEASE_LIST")
    rm -f "$unreachable_creation" "$unreachable_deletion"
    "$binpath" secrets disable "${unreachable_mount}" 2>&1 || echo "Warning: Failed to disable mount during cleanup"
    fail "Found $LEASE_COUNT phantom lease(s) - Vault should not have created leases for failed credentials"
  fi
fi
# If list command fails or returns no leases, that's the expected behavior
echo "✅ SUCCESS: No phantom leases created - Vault correctly prevented invalid credential storage"

# Cleanup
rm -f "$unreachable_creation" "$unreachable_deletion"
"$binpath" secrets disable "${unreachable_mount}" 2>&1 || echo "Warning: Failed to disable mount during cleanup"

echo "✅ All error handling tests completed successfully"
