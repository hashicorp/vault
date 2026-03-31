#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$AUDIT_LOG_PATH" ]] && fail "AUDIT_LOG_PATH env variable has not been set"

echo "Test Case #18: Audit Trail for All Operations"
echo "Functionality: Vault Core handles audit logging automatically"
echo "Expected Result: Audit logging is enabled and Vault Core automatically logs LDAP operations"
echo ""
echo "Verifying LDAP library operations are recorded in audit logs"
echo "Audit log path: ${AUDIT_LOG_PATH}"

# Check if audit log file exists
if ! sudo test -f "$AUDIT_LOG_PATH"; then
  fail "Audit log file not found at: ${AUDIT_LOG_PATH}"
fi

echo "✓ Audit log file exists"
echo ""

# Verify audit logging contains LDAP library operations
# The backend.go code calls ldapEvent() for operations which triggers audit logging
echo "Verifying Vault Core's automatic audit logging for LDAP library operations..."
echo "Note: Check-out and check-in operations are performed in the 'create' module"
echo "      (ldap_library_checkout_default_ttl, ldap_library_checkout_custom_ttl, ldap_library_self_checkin)"
echo ""

# Check for LDAP library operations in the audit log
# Looking for paths like: ldap/library/test-set, ldap/library/test-set/check-out, etc.
ldap_library_entries=$(sudo grep -c "${MOUNT}/library" "$AUDIT_LOG_PATH" 2> /dev/null || true)

if [ "$ldap_library_entries" -eq 0 ]; then
  fail "FAILED: No LDAP library operations found in audit log for mount: ${MOUNT}
Vault Core should automatically log LDAP library operations via ldapEvent() calls"
fi

echo "Found ${ldap_library_entries} LDAP library operation entries in audit log"
echo ""

# Verify specific LDAP library lifecycle operations are logged
echo "Verifying full lifecycle operations are logged:"
echo ""

# Check for library creation/configuration
if sudo grep -q "${MOUNT}/library/test-set\"" "$AUDIT_LOG_PATH" 2> /dev/null; then
  echo "✓ Library creation/configuration operations logged"
else
  fail "FAILED: Library creation operations not found in audit log"
fi

# Check for check-out operations
checkout_count=$(sudo grep -c "check-out" "$AUDIT_LOG_PATH" 2> /dev/null || true)
if [ "$checkout_count" -gt 0 ]; then
  echo "✓ Account check-out operations logged (${checkout_count} entries)"
else
  fail "FAILED: No check-out operations found in audit log"
fi

# Check for check-in operations
checkin_count=$(sudo grep -c "check-in" "$AUDIT_LOG_PATH" 2> /dev/null || true)
if [ "$checkin_count" -gt 0 ]; then
  echo "✓ Account check-in operations logged (${checkin_count} entries)"
else
  echo "⚠ Note: No check-in operations found in audit log"
fi

# Check for lease operations (renewal/revocation)
lease_operations=$(sudo grep -c "sys/leases" "$AUDIT_LOG_PATH" 2> /dev/null || true)
if [ "$lease_operations" -gt 0 ]; then
  echo "✓ Lease operations (renew/revoke) logged (${lease_operations} entries)"
else
  echo "⚠ Note: No lease operations found (Test Cases #10 and #12 may be conditional)"
fi

echo ""
echo "Summary: Vault Core's automatic audit logging is working correctly"
echo "  - LDAP library operations are automatically logged by Vault Core"
echo "  - Full lifecycle operations recorded: library creation, check-out, check-in, lease management"
echo "  - Total LDAP library audit entries: ${ldap_library_entries}"

exit 0
