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

# Create service account john for library test (test case #11)
echo "OpenLDAP: Creating service account john for library test"
SERVICE_ACCOUNT_JOHN_LDIF="service_account_john.ldif"
cat << EOF > ${SERVICE_ACCOUNT_JOHN_LDIF}
# Service Account: john
dn: uid=john,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: john
cn: john
uid: john
mail: john@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_JOHN_LDIF} 2>&1 || echo "Warning: Service account john may already exist"

# Verify the service account john was created
echo "OpenLDAP: Verifying service account john exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=john)" > /dev/null 2>&1; then
  fail "Failed to verify service account john exists in LDAP"
fi
echo "OpenLDAP: Service account john verified"

# Create service account candy for library test (test case #13)
echo "OpenLDAP: Creating service account candy for library test"
SERVICE_ACCOUNT_CANDY_LDIF="service_account_candy.ldif"
cat << EOF > ${SERVICE_ACCOUNT_CANDY_LDIF}
# Service Account: candy
dn: uid=candy,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: candy
cn: candy
uid: candy
mail: candy@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_CANDY_LDIF} 2>&1 || echo "Warning: Service account candy may already exist"

# Verify the service account candy was created
echo "OpenLDAP: Verifying service account candy exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=candy)" > /dev/null 2>&1; then
  fail "Failed to verify service account candy exists in LDAP"
fi
echo "OpenLDAP: Service account candy verified"

# Create service account sam for library test (test case #16, #17)
echo "OpenLDAP: Creating service account sam for library test"
SERVICE_ACCOUNT_SAM_LDIF="service_account_sam.ldif"
cat << EOF > ${SERVICE_ACCOUNT_SAM_LDIF}
# Service Account: sam
dn: uid=sam,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: sam
cn: sam
uid: sam
mail: sam@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_SAM_LDIF} 2>&1 || echo "Warning: Service account sam may already exist"

# Verify the service account sam was created
echo "OpenLDAP: Verifying service account sam exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=sam)" > /dev/null 2>&1; then
  fail "Failed to verify service account sam exists in LDAP"
fi
echo "OpenLDAP: Service account sam verified"

# Create service account alex for library test (test case #16, #17)
echo "OpenLDAP: Creating service account alex for library test"
SERVICE_ACCOUNT_ALEX_LDIF="service_account_alex.ldif"
cat << EOF > ${SERVICE_ACCOUNT_ALEX_LDIF}
# Service Account: alex
dn: uid=alex,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: alex
cn: alex
uid: alex
mail: alex@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_ALEX_LDIF} 2>&1 || echo "Warning: Service account alex may already exist"

# Verify the service account alex was created
echo "OpenLDAP: Verifying service account alex exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=alex)" > /dev/null 2>&1; then
  fail "Failed to verify service account alex exists in LDAP"
fi
echo "OpenLDAP: Service account alex verified"

# Create service account pat for library test
echo "OpenLDAP: Creating service account pat for library test"
SERVICE_ACCOUNT_PAT_LDIF="service_account_pat.ldif"
cat << EOF > ${SERVICE_ACCOUNT_PAT_LDIF}
# Service Account: pat
dn: uid=pat,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: pat
cn: pat
uid: pat
mail: pat@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_PAT_LDIF} 2>&1 || echo "Warning: Service account pat may already exist"

# Verify the service account pat was created
echo "OpenLDAP: Verifying service account pat exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=pat)" > /dev/null 2>&1; then
  fail "Failed to verify service account pat exists in LDAP"
fi
echo "OpenLDAP: Service account pat verified"

# Create service account kim for library test
echo "OpenLDAP: Creating service account kim for library test"
SERVICE_ACCOUNT_KIM_LDIF="service_account_kim.ldif"
cat << EOF > ${SERVICE_ACCOUNT_KIM_LDIF}
# Service Account: kim
dn: uid=kim,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: kim
cn: kim
uid: kim
mail: kim@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_KIM_LDIF} 2>&1 || echo "Warning: Service account kim may already exist"

# Verify the service account kim was created
echo "OpenLDAP: Verifying service account kim exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=kim)" > /dev/null 2>&1; then
  fail "Failed to verify service account kim exists in LDAP"
fi
echo "OpenLDAP: Service account kim verified"

# Create service account robin for library test (test case #9)
echo "OpenLDAP: Creating service account robin for library test"
SERVICE_ACCOUNT_ROBIN_LDIF="service_account_robin.ldif"
cat << EOF > ${SERVICE_ACCOUNT_ROBIN_LDIF}
# Service Account: robin
dn: uid=robin,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: robin
cn: robin
uid: robin
mail: robin@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_ROBIN_LDIF} 2>&1 || echo "Warning: Service account robin may already exist"

# Verify the service account robin was created
echo "OpenLDAP: Verifying service account robin exists"
if ! ldapsearch -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -b "ou=users,dc=${LDAP_USERNAME},dc=com" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" "(uid=robin)" > /dev/null 2>&1; then
  fail "Failed to verify service account robin exists in LDAP"
fi
echo "OpenLDAP: Service account robin verified"

# Create service account alice for enforcement test (test case #19)
echo "OpenLDAP: Creating service account alice for enforcement test"
SERVICE_ACCOUNT_ALICE_LDIF="service_account_alice.ldif"
cat << EOF > ${SERVICE_ACCOUNT_ALICE_LDIF}
# Service Account: alice
dn: uid=alice,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: alice
cn: alice
uid: alice
mail: alice@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_ALICE_LDIF} 2>&1 || echo "Warning: Service account alice may already exist"

# Create service account carol for enforcement test (test case #19)
echo "OpenLDAP: Creating service account carol for enforcement test"
SERVICE_ACCOUNT_CAROL_LDIF="service_account_carol.ldif"
cat << EOF > ${SERVICE_ACCOUNT_CAROL_LDIF}
# Service Account: carol
dn: uid=carol,ou=users,dc=${LDAP_USERNAME},dc=com
objectClass: inetOrgPerson
sn: carol
cn: carol
uid: carol
mail: carol@example.com
userPassword: ${LDAP_ADMIN_PW}
EOF
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${SERVICE_ACCOUNT_CAROL_LDIF} 2>&1 || echo "Warning: Service account carol may already exist"

echo "Vault: Configuring LDAP secrets engine at ${MOUNT}/config"
"$binpath" write "${MOUNT}/config" \
  url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
  userattr="uid" \
  insecure_tls=true
