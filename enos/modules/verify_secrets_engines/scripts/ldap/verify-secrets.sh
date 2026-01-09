#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
        echo "$1" 1>&2
        exit 1
}

# Function to generate LDIF files with DN
generate_ldif() {
        local filename="$1"
        local content="$2"
        local ou="${3:-users}"

         cat << EOF > "${filename}"
dn: cn={{.Username}},ou=${ou},dc=${LDAP_USERNAME},dc=com
${content}
EOF
}

# Validate required environment variables
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

# Set LDAP bind DN, password, and server URL for Vault's LDAP secrets engine
echo "Vault: Configuring LDAP secrets engine"
"$binpath" write ldap/config \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  url="ldap://${LDAP_SERVER}:${LDAP_PORT}"

# Create the LDIF template that Vault will use to generate dynamic LDAP users
echo "Vault: Creating creation.ldif for dynamic LDAP role"
CREATION_LDIF="creation.ldif"
generate_ldif "${CREATION_LDIF}" "objectClass: person
objectClass: top
cn: {{.Username}}
sn: {{.Password | utf16le | base64}}
userPassword: {{.Password}}" "users"

# Create the LDIF template used by Vault to delete dynamic LDAP user entries
echo "Vault: Creating deletion.ldif for dynamic LDAP role"
DELETION_LDIF="deletion.ldif"
generate_ldif "${DELETION_LDIF}" "changetype: delete" "users"

# Create the LDIF template used by Vault to roll back user creation if an error occurs
echo "Vault: Creating rollback.ldif for dynamic LDAP role"
ROLLBACK_LDIF="rollback.ldif"
generate_ldif "${ROLLBACK_LDIF}" "changetype: delete" "users"

# Create the dynamic LDAP role in Vault using the LDIF templates
echo "Vault: Creating dynamic LDAP role dynamic-role"
"$binpath" write ldap/role/dynamic-role \
  creation_ldif=@${CREATION_LDIF} \
  deletion_ldif=@${DELETION_LDIF} \
  rollback_ldif=@${ROLLBACK_LDIF} \
  default_ttl="${DEFAULT_TTL}" \
  max_ttl="${MAX_TTL}"

# Generating dynamic LDAP credentials from role dynamic-role
echo "Vault: Generating dynamic LDAP credentials from role dynamic-role"

DYNAMIC_CREDS=$("$binpath" read ldap/creds/dynamic-role)

DYNAMIC_PASSWORD=$(echo "$DYNAMIC_CREDS" | jq -r '.data.password')
DYN_DN=$(echo "$DYNAMIC_CREDS" | jq -r '.data.distinguished_names[0]')

if [[ -z "$DYN_DN" || "$DYN_DN" == "null" ]]; then
        fail "Vault did not return a distinguished_name for dynamic LDAP user"
fi

#Verify Dynamic Credentials
if ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "${DYN_DN}" -w "${DYNAMIC_PASSWORD}" > /dev/null 2>&1; then
        echo "LDAP dynamic credentials valid"
else
        fail "Error: LDAP dynamic credentials validation failed"
fi

# Attempt to use expired credentials — verify LDAP authentication fails.
CREDENTIAL_TTL_BUFFER=${CREDENTIAL_TTL_BUFFER:-80}
sleep "$CREDENTIAL_TTL_BUFFER"

if ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "$DYN_DN" -w "$DYNAMIC_PASSWORD" &> /dev/null; then
        fail "Error: Dynamic credentials still valid — TTL did NOT expire"
else
        echo "Dynamic credentials expired as expected"
fi

#Testing invalid dynamic LDAP role configuration by giving an invalid OU.

INVALID_CREATION_LDIF="invalid_creation.ldif"
INVALID_DELETION_LDIF="invalid_deletion.ldif"

generate_ldif "${INVALID_CREATION_LDIF}" "objectClass: person
objectClass: top
cn: {{.Username}}
sn: {{.Password}}
userPassword: {{.Password}}" "invalid_ou"

generate_ldif "${INVALID_DELETION_LDIF}" "changetype: delete" "invalid_ou"

echo "Vault: Attempting to create invalid-role (should fail)..."

"$binpath" write ldap/role/invalid-role \
  creation_ldif=@${INVALID_CREATION_LDIF} \
  deletion_ldif=@${INVALID_DELETION_LDIF} \
  rollback_ldif=@${INVALID_DELETION_LDIF} \
  default_ttl="${DEFAULT_TTL}" \
  max_ttl="${MAX_TTL}" || true

echo "Attempting to read creds from invalid-role"

INVALID_CREDS_OUTPUT=$(
        "$binpath" read ldap/creds/invalid-role 2>&1 || true
)
echo "$INVALID_CREDS_OUTPUT"

#check for error message indicating invalid DN/OU.
if echo "$INVALID_CREDS_OUTPUT" | grep -qi "No Such Object"; then
        echo "SUCCESS: Vault failed dynamic credential creation due to invalid OU/DN."
        "$binpath" delete ldap/role/invalid-role
else
        fail "ERROR: Vault did NOT fail when invalid DN/OU was used!"
fi
