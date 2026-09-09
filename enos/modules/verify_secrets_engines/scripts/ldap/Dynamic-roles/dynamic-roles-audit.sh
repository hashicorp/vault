#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

AUDIT_LOG="${VAULT_AUDIT_LOG:-/var/log/vault/vault_audit.log}"
ROLE_NAME="${ROLE_NAME:-dynamic-role}"

echo "Test: Audit Trail Verification for Dynamic Role Operations"
echo "Audit log path: $AUDIT_LOG"
echo "Role name: $ROLE_NAME"

if ! sudo test -f "$AUDIT_LOG"; then
  fail "Audit log file not found at $AUDIT_LOG"
fi

echo "Test 1: Checking if role creation was audited..."

if ! sudo grep -q "\"path\":\"${MOUNT}/role/${ROLE_NAME}\"" "$AUDIT_LOG"; then
  fail "ERROR: Role path not found in audit log"
fi

if sudo grep "\"path\":\"${MOUNT}/role/${ROLE_NAME}\"" "$AUDIT_LOG" | sudo grep -q "\"operation\":\"create\""; then
  echo "SUCCESS: Role creation activity found in audit log"
elif sudo grep "\"path\":\"${MOUNT}/role/${ROLE_NAME}\"" "$AUDIT_LOG" | sudo grep -q "\"operation\":\"update\""; then
  echo "SUCCESS: Role creation activity found in audit log (as update)"
else
  fail "ERROR: Role creation/update operation not found in audit log"
fi

echo "Test 2: Checking if role update was audited..."

role_entries=$(sudo grep "\"path\":\"${MOUNT}/role/${ROLE_NAME}\"" "$AUDIT_LOG")

if echo "$role_entries" | sudo grep -q "\"default_ttl\""; then
  echo "SUCCESS: Role update activity (default_ttl modification) found in audit log."
  echo "$role_entries"
elif echo "$role_entries" | sudo grep -q "\"max_ttl\""; then
  echo "SUCCESS: Role update activity (max_ttl modification) found in audit log"
else
  echo "$role_entries"
  fail "WARNING: Role update activity not clearly identified in audit log"
fi

echo "Test 3: Checking if role read was audited..."

if echo "$role_entries" | sudo grep -q "\"operation\":\"read\""; then
  echo "SUCCESS: Role read activity found in audit log"
else
  fail "ARNING: Role read activity not found in audit log"
fi
echo "All audit tests passed successfully."
