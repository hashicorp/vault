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

# Verifying LDAP Server Configs
LDAP_UID=$(ldapsearch -x -LLL -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=${LDAP_USERNAME})" 2>/dev/null)
[[ -z "$LDAP_UID" ]] && fail "Could not search ldap server for uid: ${LDAP_USERNAME}"

# Authenticate Using Vault LDAP login
VAULT_LDAP_LOGIN=$("$binpath" login -method=${MOUNT} username=${LDAP_USERNAME} password=${LDAP_ADMIN_PW})

# Verifying Vault LDAP Login Token
VAULT_LDAP_TOKEN=$(echo $VAULT_LDAP_LOGIN | jq -r ".auth.client_token")
[[ -z "$VAULT_LDAP_TOKEN" ]] && fail "Vault LDAP could not log in correctly: ${VAULT_LDAP_TOKEN}"

# Verifying Vault LDAP Policies
VAULT_POLICY_COUNT=$(echo $VAULT_LDAP_LOGIN | jq -r ".auth.policies | length")
[[ -z "$VAULT_POLICY_COUNT" ]] && fail "Vault LDAP number of policies does not look correct: ${VAULT_POLICY_COUNT}"

echo "${VAULT_LDAP_LOGIN}"
