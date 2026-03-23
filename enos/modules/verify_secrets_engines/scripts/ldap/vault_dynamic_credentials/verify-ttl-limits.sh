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
[[ -z "$DEFAULT_TTL" ]] && fail "DEFAULT_TTL env variable has not been set"
[[ -z "$MAX_TTL" ]] && fail "MAX_TTL env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo "Test: Default TTL Verification"
if ! creds=$("$binpath" read "${MOUNT}/creds/dynamic-role" 2>&1); then
  fail "Failed to generate credentials: ${creds}"
fi
initial_duration=$(jq -r '.lease_duration' <<< "$creds")
dn=$(jq -r '.data.distinguished_names[0]' <<< "$creds")
password=$(jq -r '.data.password' <<< "$creds")

[[ -z "$dn" || "$dn" == "null" ]] && fail "No distinguished_name found"
[[ -z "$password" || "$password" == "null" ]] && fail "No password found"

echo "Verifying credentials work with ldapwhoami"
if ! ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "$dn" -w "$password" > /dev/null 2>&1; then
  fail "LDAP authentication failed - credentials don't work"
fi
echo "✅ LDAP credentials verified working"

echo "Configured default_ttl: ${DEFAULT_TTL}s"
echo "Initial lease duration: ${initial_duration}s"

if [[ "$initial_duration" -eq "$DEFAULT_TTL" ]]; then
  echo "✅ SUCCESS: Initial lease duration matches default_ttl"
else
  fail "Initial lease duration ($initial_duration) does not match default_ttl ($DEFAULT_TTL)"
fi

echo "Test: Max TTL Enforcement"
if ! creds=$("$binpath" read "${MOUNT}/creds/dynamic-role" 2>&1); then
  fail "Failed to generate credentials: ${creds}"
fi
lease_id=$(jq -r '.lease_id' <<< "$creds")

[[ -z "$lease_id" || "$lease_id" == "null" ]] && fail "No lease_id found"

echo "Renewing credential to max_ttl limit (${MAX_TTL}s)"
if ! renewed=$("$binpath" write sys/leases/renew lease_id="$lease_id" increment="${MAX_TTL}s" 2>&1); then
  fail "Failed to renew credential: ${renewed}"
fi
echo "Renewal output:"
jq '.' <<< "$renewed"
renewed_ttl=$(jq -r '.lease_duration' <<< "$renewed")

if [[ "$renewed_ttl" -le "$MAX_TTL" ]]; then
  echo "✅ SUCCESS: Renewed TTL ($renewed_ttl) respects max_ttl ($MAX_TTL)"
else
  fail "Renewed TTL ($renewed_ttl) exceeds max_ttl ($MAX_TTL)"
fi

echo "Attempting to renew beyond max_ttl"
if ! beyond_max_result=$("$binpath" write sys/leases/renew lease_id="$lease_id" increment="$((MAX_TTL * 2))s" 2>&1); then
  fail "Failed to renew credential beyond max_ttl: ${beyond_max_result}"
fi
echo "Renewal beyond max_ttl output:"
jq '.' <<< "$beyond_max_result"
capped_ttl=$(jq -r '.lease_duration' <<< "$beyond_max_result")

if [[ "$capped_ttl" -le "$MAX_TTL" ]]; then
  echo "✅ SUCCESS: TTL capped at ${capped_ttl}s when attempting to renew beyond max_ttl (${MAX_TTL}s)"
else
  fail "TTL ($capped_ttl) was not capped to max_ttl ($MAX_TTL)"
fi
