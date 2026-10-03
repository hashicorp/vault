#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "ERROR: $1" 1>&2
  exit 1
}

[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER not set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT not set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME not set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW not set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR not set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR not set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN not set"

binpath="${VAULT_INSTALL_DIR}/vault"
[[ -x "$binpath" ]] || fail "Vault binary not found"

export VAULT_FORMAT=json

LDAP_URL="ldap://${LDAP_SERVER}:${LDAP_PORT}"
BINDDN="cn=admin,dc=${LDAP_USERNAME},dc=com"
BAD_BINDDN="cn=nonexistent-admin,dc=${LDAP_USERNAME},dc=com"

echo "Test: LDAP credential mismatch behavior"

# Baseline validation
ldapwhoami -x -H "$LDAP_URL" -D "$BINDDN" \
  -w "$LDAP_ADMIN_PW" > /dev/null 2>&1 \
  || fail "Baseline LDAP credentials do not work"

# Poison Vault config (wrong binddn + wrong password)
"$binpath" write ldap/config \
  binddn="$BAD_BINDDN" \
  bindpass="intentionally-wrong-password" \
  url="$LDAP_URL" > /dev/null

# Attempt rotation (observational)
set +e
ROTATION_OUTPUT=$("$binpath" write -f ldap/rotate-root 2>&1)
ROTATION_EXIT=$?
set -e

echo "$ROTATION_OUTPUT"

if [[ $ROTATION_EXIT -ne 0 ]]; then
  echo "Rotation failed as expected "
else
  echo "Rotation succeeded, unexpected"
fi

# Restore correct config
"$binpath" write ldap/config \
  binddn="$BINDDN" \
  bindpass="$LDAP_ADMIN_PW" \
  url="$LDAP_URL" > /dev/null

# Post-recovery validation
"$binpath" read ldap/config > /dev/null 2>&1 \
  || fail "Vault failed to reconnect after recovery"

ldapwhoami -x -H "$LDAP_URL" -D "$BINDDN" \
  -w "$LDAP_ADMIN_PW" > /dev/null 2>&1 \
  || fail "LDAP authentication failed after recovery"

echo "Test complete"
exit 0
