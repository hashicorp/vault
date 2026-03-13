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
[[ -z "$LDAP_DOMAIN" ]] && fail "LDAP_DOMAIN env variable has not been set"

echo "OpenLDAP: Checking for OpenLDAP Server Connection: ${LDAP_SERVER}:${LDAP_PORT}"
# Wait for LDAP server to be ready
sleep 10

# Extract domain components from LDAP_DOMAIN (e.g., "enos.com" -> "dc=enos,dc=com")
IFS='.' read -ra DOMAIN_PARTS <<< "$LDAP_DOMAIN"
DOMAIN_DN=""
for part in "${DOMAIN_PARTS[@]}"; do
  if [[ -n "$DOMAIN_DN" ]]; then
    DOMAIN_DN="${DOMAIN_DN},dc=${part}"
  else
    DOMAIN_DN="dc=${part}"
  fi
done

echo "OpenLDAP: Using domain DN: ${DOMAIN_DN}"
echo "OpenLDAP: Testing connection with admin credentials"

# Test connection
ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "${DOMAIN_DN}" -D "cn=admin,${DOMAIN_DN}" -w "${LDAP_ADMIN_PW}" -s base

echo "OpenLDAP: Creating organizational units"
# Creating Users and Groups Org Units LDIF file
OU_LDIF="ou.ldif"
cat << EOF > ${OU_LDIF}
dn: ou=users,${DOMAIN_DN}
objectClass: organizationalUnit
ou: users

dn: ou=groups,${DOMAIN_DN}
objectClass: organizationalUnit
ou: groups
EOF

ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,${DOMAIN_DN}" -w "${LDAP_ADMIN_PW}" -f ${OU_LDIF} || echo "OUs may already exist"

echo "OpenLDAP: Creating test users"
USER_LDIF="users.ldif"
cat << EOF > ${USER_LDIF}
# User: enos
dn: uid=enos,ou=users,${DOMAIN_DN}
objectClass: inetOrgPerson
sn: enos
cn: enos user
uid: enos
userPassword: ${LDAP_ADMIN_PW}

# Static-role test user (for LDAP verification tests)
dn: uid=vault-static-user,ou=users,${DOMAIN_DN}
objectClass: inetOrgPerson
sn: vault-static-user
cn: Vault Static User
uid: vault-static-user
userPassword: ${LDAP_ADMIN_PW}

# Service accounts for library tests
dn: uid=svc-account-1,ou=users,${DOMAIN_DN}
objectClass: inetOrgPerson
sn: svc-account-1
cn: Service Account 1
uid: svc-account-1
userPassword: ${LDAP_ADMIN_PW}

dn: uid=svc-account-2,ou=users,${DOMAIN_DN}
objectClass: inetOrgPerson
sn: svc-account-2
cn: Service Account 2
uid: svc-account-2
userPassword: ${LDAP_ADMIN_PW}

dn: uid=svc-delete,ou=users,${DOMAIN_DN}
objectClass: inetOrgPerson
sn: svc-delete
cn: Service Account Delete
uid: svc-delete
userPassword: ${LDAP_ADMIN_PW}
EOF

ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,${DOMAIN_DN}" -w "${LDAP_ADMIN_PW}" -f ${USER_LDIF} || echo "Users may already exist"

echo "LDAP population completed successfully."
