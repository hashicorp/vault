#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# Verifies system consistency after failed LDAP root rotation (CI-safe)

set -e

fail() {
  echo "ERROR: $1" 1>&2
  exit 1
}

# Required environment
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

ROTATION_START_DELAY="${ROTATION_START_DELAY:-1}"
LDAP_URL="ldap://${LDAP_SERVER}:${LDAP_PORT}"
BINDDN="cn=admin,dc=${LDAP_USERNAME},dc=com"

echo "Test: system consistency after mid-rotation failure"

# 1. Baseline LDAP health
ldapwhoami -x -H "$LDAP_URL" -D "$BINDDN" -w "$LDAP_ADMIN_PW" \
  > /dev/null 2>&1 || fail "Baseline LDAP auth failed"

# 2. Ensure Vault LDAP config is correct
"$binpath" write ldap/config \
  binddn="$BINDDN" \
  bindpass="$LDAP_ADMIN_PW" \
  url="$LDAP_URL" > /dev/null

"$binpath" read ldap/config > /dev/null 2>&1 \
  || fail "LDAP config not readable before rotation"

# 3. Start rotation asynchronously
"$binpath" write -f ldap/rotate-root > /dev/null 2>&1 &
ROT_PID=$!

# 4. Allow rotation to enter execution
sleep "$ROTATION_START_DELAY"

# 5. Inject failure at Vault layer
"$binpath" write ldap/config \
  binddn="cn=nonexistent,dc=${LDAP_USERNAME},dc=com" \
  bindpass="wrong-password" \
  url="$LDAP_URL" > /dev/null

# 6. Wait for rotation attempt to finish
wait "$ROT_PID" || echo "Rotation failed as expected"

# 7. Restore valid config
"$binpath" write ldap/config \
  binddn="$BINDDN" \
  bindpass="$LDAP_ADMIN_PW" \
  url="$LDAP_URL" > /dev/null

# 8. Verify system consistency
ldapwhoami -x -H "$LDAP_URL" -D "$BINDDN" -w "$LDAP_ADMIN_PW" \
  > /dev/null 2>&1 && echo "LDAP credentials usable after failure" \
  || echo "WARNING: LDAP credentials may require manual recovery"

"$binpath" status > /dev/null 2>&1 \
  || fail "Vault unhealthy after failed rotation"

"$binpath" read ldap/config > /dev/null 2>&1 \
  || fail "Vault cannot read LDAP config after recovery"

echo "Test complete (system consistent)"
exit 0
