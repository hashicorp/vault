#!/usr/bin/env bash
# Copyright IBM Corp. 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$LDAP_BASE_DN" ]] && fail "LDAP_BASE_DN env variable has not been set"

echo "OpenLDAP: Checking for OpenLDAP Server Connection: ${LDAP_SERVER}:${LDAP_PORT}"
echo "OpenLDAP: Using base DN: ${LDAP_BASE_DN}"
echo "OpenLDAP: Testing connection with admin credentials"

echo "OpenLDAP: Creating organizational units"
# Creating Users and Groups Org Units LDIF file
OU_LDIF="ou.ldif"
cat << EOF > "${OU_LDIF}"
dn: ou=users,${LDAP_BASE_DN}
objectClass: organizationalUnit
ou: users

dn: ou=groups,${LDAP_BASE_DN}
objectClass: organizationalUnit
ou: groups
EOF

ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,${LDAP_BASE_DN}" -w "${LDAP_ADMIN_PW}" -f "${OU_LDIF}" || echo "OUs may already exist"

echo "OpenLDAP: Creating test users"
USER_LDIF="users.ldif"
cat << EOF > "${USER_LDIF}"
# User: enos
dn: uid=enos,ou=users,${LDAP_BASE_DN}
objectClass: inetOrgPerson
sn: enos
cn: enos user
uid: enos
userPassword: ${LDAP_ADMIN_PW}

# Static-role test user (for LDAP verification tests)
dn: uid=vault-static-user,ou=users,${LDAP_BASE_DN}
objectClass: inetOrgPerson
sn: vault-static-user
cn: Vault Static User
uid: vault-static-user
userPassword: ${LDAP_ADMIN_PW}

# Service accounts for library tests
dn: uid=svc-account-1,ou=users,${LDAP_BASE_DN}
objectClass: inetOrgPerson
sn: svc-account-1
cn: Service Account 1
uid: svc-account-1
userPassword: ${LDAP_ADMIN_PW}

dn: uid=svc-account-2,ou=users,${LDAP_BASE_DN}
objectClass: inetOrgPerson
sn: svc-account-2
cn: Service Account 2
uid: svc-account-2
userPassword: ${LDAP_ADMIN_PW}

dn: uid=svc-delete,ou=users,${LDAP_BASE_DN}
objectClass: inetOrgPerson
sn: svc-delete
cn: Service Account Delete
uid: svc-delete
userPassword: ${LDAP_ADMIN_PW}
EOF

ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,${LDAP_BASE_DN}" -w "${LDAP_ADMIN_PW}" -f "${USER_LDIF}" || echo "Users may already exist"

echo "LDAP population completed successfully."
