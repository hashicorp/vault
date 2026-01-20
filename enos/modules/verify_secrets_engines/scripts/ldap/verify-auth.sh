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
# Authenticate Using Vault LDAP login
VAULT_LDAP_LOGIN=$("$binpath" login -format=json -method="${MOUNT}" username="${LDAP_USERNAME}" password="${LDAP_ADMIN_PW}")

# Verify that admin-policy is attached to the token
verify_admin_policy_attachment() {
  if ! echo "$VAULT_LDAP_LOGIN" \
      | jq -e '.auth.policies[] | select(. == "admin-policy")' > /dev/null; then
    fail "admin-policy is NOT attached to the LDAP login token"
  fi
}

# Verify LDAP user is in vault-admins group
verify_ldap_group_membership() {
  if ! ldapsearch -x \
      -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" \
      -D "cn=admin,dc=${LDAP_USERNAME},dc=com" \
      -w "${LDAP_ADMIN_PW}" \
      -b "ou=groups,dc=${LDAP_USERNAME},dc=com" \
      "(&(cn=vault-admins)(member=uid=${LDAP_USERNAME},ou=users,dc=${LDAP_USERNAME},dc=com))" \
      | grep -qi '^dn:'; then
    fail "User NOT in vault-admins group"
  fi
}

# Verify vault-admins mapped to admin-policy
verify_policy_mapping() {
  if ! "$binpath" read -format=json auth/"${MOUNT}"/groups/vault-admins \
       | jq -e '.data.policies[] | select(. == "admin-policy")' > /dev/null; then
    fail "vault-admins is NOT mapped to admin-policy"
  fi
}

# Verify LDAP admin login successful
verify_admin_login() {
  if ! ADMIN_LOGIN_OUTPUT=$(
    "$binpath" login -format=json \
      -method="${MOUNT}" \
      username="${LDAP_USERNAME}" \
      password="${LDAP_ADMIN_PW}"
  ); then
    fail "LDAP admin login FAILED"
  fi
  jq -e '.auth.client_token' <<< "$ADMIN_LOGIN_OUTPUT" > /dev/null \
    || fail "LDAP admin login FAILED"
}

# Verify wrong passwordlogin fails
verify_wrong_password_login_fails() {
  BAD_LOGIN_OUTPUT=$(
    "$binpath" login \
      -method="${MOUNT}" \
      username="${LDAP_USERNAME}" \
      password="wrong-password" 2>&1 || true
  )
  if ! echo "$BAD_LOGIN_OUTPUT" | grep -Eqi "(invalid|failed to bind as user)"; then
    fail "Security Failure: Login succeeded with wrong password"
  fi
}

# Deleting admin policy
cleanup_admin_policy() {
  echo "Vault: Deleting admin policy"
  "$binpath" policy delete "admin-policy" > /dev/null 2>&1 || true
}

verify_admin_policy_attachment
verify_ldap_group_membership
verify_policy_mapping
verify_admin_login
verify_wrong_password_login_fails
cleanup_admin_policy
echo "${VAULT_LDAP_LOGIN}"
