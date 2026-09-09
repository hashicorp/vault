#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"

echo "OpenLDAP: Creating password policies organizational unit LDIF file"
PWPOLICIES_OU_LDIF="pwpolicies-ou.ldif"
cat << EOF > ${PWPOLICIES_OU_LDIF}
dn: ou=pwpolicies,dc=${LDAP_USERNAME},dc=com
objectClass: organizationalUnit
ou: pwpolicies
EOF

echo "OpenLDAP: Adding password policies organizational unit"
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${PWPOLICIES_OU_LDIF}
echo "OpenLDAP: Password policies organizational unit added successfully"

echo "OpenLDAP: Creating ppolicy module LDIF file"
PASSWORD_POLICY_MODULE_LDIF="password-policy-module.ldif"
cat << EOF > ${PASSWORD_POLICY_MODULE_LDIF}
dn: cn=module{0},cn=config
changetype: modify
add: olcModuleLoad
olcModuleLoad: ppolicy
EOF

echo "OpenLDAP: Loading ppolicy module"
ldapmodify -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,cn=config" -w config -f ${PASSWORD_POLICY_MODULE_LDIF}
echo "OpenLDAP: ppolicy module loaded successfully"

echo "OpenLDAP: Creating ppolicy overlay LDIF file"
PASSWORD_POLICY_OVERLAY_LDIF="password-policy-overlay.ldif"
cat << EOF > ${PASSWORD_POLICY_OVERLAY_LDIF}
dn: olcOverlay=ppolicy,olcDatabase={1}mdb,cn=config
objectClass: olcOverlayConfig
objectClass: olcPPolicyConfig
olcOverlay: ppolicy
olcPPolicyDefault: cn=default,ou=pwpolicies,dc=${LDAP_USERNAME},dc=com
olcPPolicyHashCleartext: TRUE
olcPPolicyUseLockout: TRUE
EOF

echo "OpenLDAP: Adding ppolicy overlay"
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,cn=config" -w config -f ${PASSWORD_POLICY_OVERLAY_LDIF}
echo "OpenLDAP: ppolicy overlay added successfully"

echo "OpenLDAP: Creating default password policy LDIF file"
DEFAULT_POLICY_LDIF="default-policy.ldif"
cat << EOF > ${DEFAULT_POLICY_LDIF}
dn: cn=default,ou=pwpolicies,dc=${LDAP_USERNAME},dc=com
objectClass: top
objectClass: organizationalRole
objectClass: pwdPolicy
cn: default
pwdAttribute: userPassword
pwdMinLength: 8
pwdInHistory: 5
pwdMaxAge: 7776000
pwdMinAge: 0
pwdAllowUserChange: TRUE
pwdExpireWarning: 604800
pwdGraceAuthNLimit: 3
pwdLockout: TRUE
pwdLockoutDuration: 1800
pwdMaxFailure: 5
pwdFailureCountInterval: 300
pwdMustChange: FALSE
pwdCheckQuality: 2
pwdSafeModify: FALSE
EOF

echo "OpenLDAP: Adding default password policy"
ldapadd -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "cn=admin,dc=${LDAP_USERNAME},dc=com" -w "${LDAP_ADMIN_PW}" -f ${DEFAULT_POLICY_LDIF}
echo "OpenLDAP: Default password policy added successfully"
