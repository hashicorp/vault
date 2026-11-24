#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# Generic audit trail verification for Vault operations

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${VAULT_AUDIT_LOG}" ]] && fail "VAULT_AUDIT_LOG env variable has not been set"

export VAULT_FORMAT=json

[[ -f "$VAULT_AUDIT_LOG" ]] || fail "Audit log file does not exist at $VAULT_AUDIT_LOG"

# Verify rotate-root event
if rotate_root=$(sudo grep -E "rotate-root" "$VAULT_AUDIT_LOG" 2>&1); then
  echo "SUCCESS: Audit log contains rotate-root event"
  echo "Found entries: $rotate_root"
else
  fail "Audit log does not contain rotate-root event. Grep output: $rotate_root"
fi

# Verify constraint violation event (if expected)
if constraint_output=$(sudo grep -E "Constraint Violation" "$VAULT_AUDIT_LOG" 2>&1); then
  echo "SUCCESS: Audit log contains constraint violation event"
  echo "Found entries: $constraint_output"
else
  fail "Audit log does not contain constraint violation event. Grep output: $constraint_output"
fi
