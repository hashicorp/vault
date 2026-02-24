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

echo "Test: Automatic Cleanup When Lease Expires"
if ! creds=$("$binpath" read "${MOUNT}/creds/dynamic-role" 2>&1); then
  fail "Failed to generate credential for expiry test: ${creds}"
fi
lease_id=$(jq -r '.lease_id' <<< "$creds")
lease_duration=$(jq -r '.lease_duration' <<< "$creds")
dn=$(jq -r '.data.distinguished_names[0]' <<< "$creds")
password=$(jq -r '.data.password' <<< "$creds")
username=$(cut -d',' -f1 <<< "$dn" | cut -d'=' -f2)

[[ -z "$lease_id" || "$lease_id" == "null" ]] && fail "No lease_id found"
[[ -z "$dn" || "$dn" == "null" ]] && fail "No distinguished_name found"
[[ -z "$password" || "$password" == "null" ]] && fail "No password found"
[[ -z "$lease_duration" || "$lease_duration" == "null" ]] && fail "No lease_duration found"

echo "Verifying credentials work with ldapwhoami: $username"
if ! ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "$dn" -w "$password" > /dev/null 2>&1; then
  fail "LDAP authentication failed - credentials don't work"
fi
echo "✅ LDAP credentials verified working"

echo "Waiting for lease to expire (${lease_duration}s)..."
sleep "$lease_duration"

echo "Attempting to renew expired lease (should fail)..."
if "$binpath" write sys/leases/renew lease_id="$lease_id" increment="60s" > /dev/null 2>&1; then
  fail "Lease renewal should have failed but succeeded"
else
  echo "✅ SUCCESS: Lease renewal correctly failed with non-zero exit code - lease has expired and been revoked"
fi

echo ""
echo "Verifying credentials no longer work (automatic cleanup)"
if ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "$dn" -w "$password" > /dev/null 2>&1; then
  fail "Credentials still work after lease expiration - automatic cleanup failed"
fi
echo "✅ SUCCESS: Credentials no longer work - automatic cleanup completed"
