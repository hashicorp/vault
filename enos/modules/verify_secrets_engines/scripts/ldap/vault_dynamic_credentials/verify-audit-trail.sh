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
[[ -z "$VAULT_AUDIT_LOG" ]] && fail "VAULT_AUDIT_LOG env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# Audit log configuration - path must be explicitly provided via VAULT_AUDIT_LOG
audit_log="${VAULT_AUDIT_LOG}"
audit_poll_timeout="${AUDIT_POLL_TIMEOUT:-30}"
audit_poll_interval="${AUDIT_POLL_INTERVAL:-1}"

# Function to poll audit log for a pattern with timeout
poll_audit_log() {
  local pattern="$1"
  local timeout="$2"
  local interval="${3:-1}"
  local description="${4:-polling audit log}"
  local elapsed=0

  while [[ $elapsed -lt $timeout ]]; do
    if sudo grep -q "$pattern" "$audit_log" 2> /dev/null; then
      return 0
    fi
    sleep "$interval"
    elapsed=$((elapsed + interval))
  done
  echo "ERROR: Timed out after ${timeout}s while ${description}. Pattern searched: ${pattern}" >&2
  return 1
}

echo "Test: Audit Trail Verification for Dynamic Credentials"
echo "Audit log path: $audit_log"

echo ""
echo "=== Creating New Credential and Verifying Audit Trail ==="
if ! creds=$("$binpath" read "${MOUNT}/creds/dynamic-role" 2>&1); then
  fail "Failed to generate credential: ${creds}"
fi

echo "Credential generated successfully"
echo "Polling audit log for credential creation event (timeout: ${audit_poll_timeout}s)..."

# Check for credential creation in audit log
if ! poll_error=$(poll_audit_log "${MOUNT}/creds/dynamic-role" "$audit_poll_timeout" "$audit_poll_interval" "verifying credential creation event in audit log" 2>&1); then
  fail "New credential creation event not found in audit log. ${poll_error}"
fi
echo "✅ SUCCESS: New credential creation event found in audit log"

echo ""
echo "=== Audit Trail Verification Complete ==="
echo "Credential creation has been successfully audited in: $audit_log"
