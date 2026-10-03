#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "ERROR: $1" 1>&2
  exit 1
}

[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath="${VAULT_INSTALL_DIR}/vault"
[[ -x "$binpath" ]] || fail "Vault binary not found at $binpath"

export VAULT_FORMAT=json

AUDIT_LOG="${VAULT_AUDIT_LOG:-/var/log/vault/vault_audit.log}"
CONFIG_RESTORE_DELAY="${CONFIG_RESTORE_DELAY:-3}"

LDAP_URL="ldap://${LDAP_SERVER}:${LDAP_PORT}"
BAD_LDAP_URL="ldap://${LDAP_SERVER}:9999"
BINDDN="cn=admin,dc=${LDAP_USERNAME},dc=com"

echo "Test: invalid LDAP endpoint prevents rotation and preserves credentials"

# Baseline validation
ldapwhoami -x -H "$LDAP_URL" -D "$BINDDN" \
  -w "$LDAP_ADMIN_PW" > /dev/null 2>&1 \
  || fail "Baseline LDAP credentials do not work"

# Poison Vault with unreachable LDAP endpoint
"$binpath" write ldap/config \
  binddn="$BINDDN" \
  bindpass="$LDAP_ADMIN_PW" \
  url="$BAD_LDAP_URL" > /dev/null

# Attempt rotation
ROTATION_OUTPUT=$("$binpath" write -f ldap/rotate-root 2>&1 || true)
ROTATION_EXIT=$?

echo "$ROTATION_OUTPUT"

if [[ $ROTATION_EXIT -ne 0 ]] || echo "$ROTATION_OUTPUT" | grep -qiE "error|fail|connection|timeout"; then
  echo "Rotation failed as expected"
else
  echo "WARNING: rotation did not fail as expected"
fi

# Restore valid config
"$binpath" write ldap/config \
  binddn="$BINDDN" \
  bindpass="$LDAP_ADMIN_PW" \
  url="$LDAP_URL" > /dev/null

sleep "$CONFIG_RESTORE_DELAY"

# Verify credentials preserved
ldapwhoami -x -H "$LDAP_URL" -D "$BINDDN" \
  -w "$LDAP_ADMIN_PW" > /dev/null 2>&1 \
  || fail "LDAP credentials were unexpectedly changed"

echo "Credentials preserved"

"$binpath" read ldap/config > /dev/null 2>&1 && echo "Vault reconnected"

grep -E "rotate-root" "$AUDIT_LOG" > /dev/null 2>&1 && echo "Audit log updated"

echo "Test complete"
exit 0
