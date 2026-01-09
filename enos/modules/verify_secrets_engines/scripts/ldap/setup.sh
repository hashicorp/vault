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

echo "OpenLDAP: Checking for OpenLDAP Server Connection: ${LDAP_SERVER}:${LDAP_PORT}"
ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}"

# Creating Users Org Unit LDIF file and adding users organizational unit
echo "OpenLDAP: Creating Users Org Unit LDIF file and adding users organizational unit"
GROUP_LDIF="group.ldif"
cat << EOF > ${GROUP_LDIF}
dn: ou=users,dc=$LDAP_USERNAME,dc=com
objectClass: organizationalUnit
ou: users

dn: ou=groups,dc=$LDAP_USERNAME,dc=com
objectClass: organizationalUnit
ou: groups
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${GROUP_LDIF}

ADMIN_LDIF="admin.ldif"
cat << EOF > ${ADMIN_LDIF}
dn: cn=admin,dc=$LDAP_USERNAME,dc=com
objectClass: simpleSecurityObject
objectClass: organizationalRole
cn: admin
description: LDAP administrator
userPassword: $LDAP_ADMIN_PW
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${ADMIN_LDIF}

echo "OpenLDAP: Creating User LDIF file and adding user to LDAP"
USER_LDIF="user.ldif"
cat << EOF > ${USER_LDIF}
# User: enos
dn: uid=$LDAP_USERNAME,ou=users,dc=$LDAP_USERNAME,dc=com
objectClass: inetOrgPerson
sn: $LDAP_USERNAME
cn: $LDAP_USERNAME user
uid: $LDAP_USERNAME
userPassword: $LDAP_ADMIN_PW

# Group: devs
dn: cn=devs,ou=groups,dc=$LDAP_USERNAME,dc=com
objectClass: groupOfNames
cn: devs
member: uid=$LDAP_USERNAME,ou=users,dc=$LDAP_USERNAME,dc=com
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${USER_LDIF}
