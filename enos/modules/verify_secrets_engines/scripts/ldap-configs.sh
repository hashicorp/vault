#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

install_ldap_tools() {
  echo "Installing OpenLDAP client tools..."
  if [[ "$OSTYPE" == "linux"* ]]; then
    if [ -x "$(command -v apt-get)" ]; then
      sudo apt-get update
      sudo apt-get install -y ldap-utils
    elif [ -x "$(command -v yum)" ]; then
      sudo yum install -y openldap-clients
    elif [ -x "$(command -v zypper)" ]; then
      sudo zypper install -y openldap2-client
    else
      echo "Unsupported Linux package manager."
      exit 1
    fi
  else
    echo "Unsupported OS: $OSTYPE"
    exit 1
  fi
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$LDAP_HOST" ]] && fail "LDAP_HOST env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# Installing LDAP tools
install_ldap_tools > /dev/null

echo "OpenLDAP: Checking for OpenLDAP Server Connection: ${LDAP_HOST}:${LDAP_PORT}"
ldapsearch -x -H "ldap://${LDAP_HOST}:${LDAP_PORT}" -b "dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}"

# Creating Users Org Unit LDIF file and adding users organizational unit
echo "OpenLDAP: Creating Users Org Unit LDIF file and adding users organizational unit"
GROUP_LDIF="group.ldif"
cat << EOF > ${GROUP_LDIF}
dn: ou=users,dc=$LDAP_USERNAME,dc=com
objectClass: organizationalUnit
ou: users
EOF
ldapadd -x -H "ldap://${LDAP_HOST}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${GROUP_LDIF}

echo "OpenLDAP: Creating User LDIF file and adding user to LDAP"
USER_LDIF="user.ldif"
cat << EOF > ${USER_LDIF}
dn: uid=$LDAP_USERNAME,ou=users,dc=$LDAP_USERNAME,dc=com
objectClass: inetOrgPerson
sn: user
cn: $LDAP_USERNAME user
uid: $LDAP_USERNAME
userPassword: $LDAP_ADMIN_PW
EOF
ldapadd -x -H "ldap://${LDAP_HOST}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${USER_LDIF}

echo "Vault: Creating ldap auth and creating auth/ldap/config route"
"$binpath" auth enable "${MOUNT}" > /dev/null 2>&1 || echo "Warning: Vault ldap auth already enabled"
"$binpath" write "auth/${MOUNT}/config" \
  url="ldap://test_${LDAP_HOST}:${LDAP_PORT}" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
  userattr="uid" \
  groupdn="ou=groups,dc=${LDAP_USERNAME},dc=com" \
  groupfilter="(&(objectClass=groupOfNames)(member={{.UserDN}}))" \
  groupattr="cn" \
  insecure_tls=true

echo "Vault: Updating ldap auth and creating auth/ldap/config route"
"$binpath" auth enable "${MOUNT}" > /dev/null 2>&1 || echo "Warning: Vault ldap auth already enabled"
"$binpath" write "auth/${MOUNT}/config" \
  url="ldap://${LDAP_HOST}:${LDAP_PORT}" \
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
"$binpath" policy write reader "${VAULT_LDAP_POLICY}"
"$binpath" write "auth/${MOUNT}/users/${LDAP_USERNAME}" policies="reader"
