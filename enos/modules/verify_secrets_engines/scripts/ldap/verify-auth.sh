#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

# Validate required environment variables
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

echo "Vault: Creating ldap auth and creating auth/ldap/config route"
"$binpath" auth enable "${MOUNT}" > /dev/null 2>&1 || echo "Warning: Vault ldap auth already enabled"
"$binpath" write "auth/${MOUNT}/config" \
  url="ldap://test_${LDAP_SERVER}:${LDAP_PORT}" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
  userattr="uid" \
  groupdn="ou=groups,dc=${LDAP_USERNAME},dc=com" \
  groupfilter="(&(objectClass=groupOfNames)(member={{.UserDN}}))" \
  groupattr="cn" \
  insecure_tls=true

echo "Vault: Updating ldap auth and creating auth/ldap/config route"
"$binpath" write "auth/${MOUNT}/config" \
  url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
  userattr="uid" \
  groupdn="ou=groups,dc=${LDAP_USERNAME},dc=com" \
  groupfilter="(&(objectClass=groupOfNames)(member={{.UserDN}}))" \
  groupattr="cn" \
  insecure_tls=true

echo "Vault: Creating Vault Policy for LDAP and assigning user to policy"
VAULT_LDAP_POLICY="ldap_reader.hcl"
cat << EOF > ${VAULT_LDAP_POLICY}
path "secret/data/*" {
  capabilities = ["read", "list"]
}
EOF
LDAP_READER_POLICY="reader-policy"
"$binpath" policy write ${LDAP_READER_POLICY} "${VAULT_LDAP_POLICY}"
"$binpath" write "auth/${MOUNT}/users/${LDAP_USERNAME}" policies="${LDAP_READER_POLICY}"

echo "Vault: Creating Vault Policy for LDAP DEV and assigning user to policy"
VAULT_LDAP_DEV_POLICY="ldap_dev.hcl"
cat << EOF > ${VAULT_LDAP_DEV_POLICY}
path "secret/data/dev/*" {
  capabilities = ["read", "list"]
}
EOF
LDAP_DEV_POLICY="dev-policy"
"$binpath" policy write ${LDAP_DEV_POLICY} "${VAULT_LDAP_DEV_POLICY}"
"$binpath" write "auth/${MOUNT}/groups/devs" policies="${LDAP_DEV_POLICY}"

echo "SUCCESS: LDAP auth engine configured and verified"
