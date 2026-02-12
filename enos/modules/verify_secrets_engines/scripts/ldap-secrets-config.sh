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

# Wait for LDAP server to be ready
echo "OpenLDAP: Waiting for LDAP server to be ready at ${LDAP_SERVER}:${LDAP_PORT}"
for i in {1..60}; do
  if ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -s base > /dev/null 2>&1; then
    echo "OpenLDAP: Server is ready"
    break
  fi
  if [ "$i" -eq 60 ]; then
    fail "LDAP server did not become ready after 60 attempts (2 minutes)"
  fi
  sleep 2
done

# Create service account in LDAP for library test
echo "OpenLDAP: Creating service account fizz for library test"
SERVICE_ACCOUNT_LDIF="service_account.ldif"
cat << EOF > ${SERVICE_ACCOUNT_LDIF}
# Service Account: fizz (uid matches userattr configuration)
dn: uid=fizz,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: fizz
cn: fizz
uid: fizz
mail: fizz@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_LDIF} 2>&1 || echo "Warning: Service account may already exist"

# Verify the service account was created
echo "OpenLDAP: Verifying service account fizz exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=fizz)" > /dev/null 2>&1; then
  fail "Failed to verify service account fizz exists in LDAP"
fi
echo "OpenLDAP: Service account fizz verified"

# Create service account buzz in LDAP for library update test (test case #3)
echo "OpenLDAP: Creating service account buzz for library update test"
SERVICE_ACCOUNT_BUZZ_LDIF="service_account_buzz.ldif"
cat << EOF > ${SERVICE_ACCOUNT_BUZZ_LDIF}
# Service Account: buzz (uid matches userattr configuration)
dn: uid=buzz,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: buzz
cn: buzz
uid: buzz
mail: buzz@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_BUZZ_LDIF} 2>&1 || echo "Warning: Service account buzz may already exist"

# Verify the service account buzz was created
echo "OpenLDAP: Verifying service account buzz exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=buzz)" > /dev/null 2>&1; then
  fail "Failed to verify service account buzz exists in LDAP"
fi
echo "OpenLDAP: Service account buzz verified"

echo "Vault: Configuring LDAP secrets engine at ${MOUNT}/config"
"$binpath" write "${MOUNT}/config" \
  url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
  userattr="uid" \
  insecure_tls=true
