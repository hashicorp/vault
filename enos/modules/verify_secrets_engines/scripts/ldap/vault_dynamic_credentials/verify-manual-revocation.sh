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

echo "Test: Manual Credential Revocation"
if ! creds=$("$binpath" read "${MOUNT}/creds/dynamic-role" 2>&1); then
  fail "Failed to generate credential: ${creds}"
fi
lease_id=$(jq -r '.lease_id' <<< "$creds")
dn=$(jq -r '.data.distinguished_names[0]' <<< "$creds")
password=$(jq -r '.data.password' <<< "$creds")
username=$(cut -d',' -f1 <<< "$dn" | cut -d'=' -f2)

[[ -z "$lease_id" || "$lease_id" == "null" ]] && fail "No lease_id found"
[[ -z "$dn" || "$dn" == "null" ]] && fail "No distinguished_name found"
[[ -z "$password" || "$password" == "null" ]] && fail "No password found"

echo "Verifying credentials work with ldapwhoami: $username"
if ! ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "$dn" -w "$password" > /dev/null 2>&1; then
  fail "LDAP authentication failed - credentials don't work"
fi
echo "✅ LDAP credentials verified working"

echo "Revoking credential"
if ! revoke_output=$("$binpath" write sys/leases/revoke lease_id="$lease_id" 2>&1); then
  fail "Failed to revoke credential: ${revoke_output}"
fi

echo "Verifying lease was revoked by attempting renewal"
if "$binpath" write sys/leases/renew lease_id="$lease_id" > /dev/null 2>&1; then
  fail "Lease renewal succeeded - credential was not revoked"
else
  echo "✅ SUCCESS: Lease renewal failed with non-zero exit code - credential was successfully revoked"
fi
