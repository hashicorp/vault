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

# Function to create a password policy in Vault
create_password_policy() {
    local policy_name="$1"
    local length="$2"
    local lowercase_min="$3"
    local uppercase_min="$4"
    local digit_min="$5"
    local special_chars="$6"
    local special_min="$7"

    echo "Vault: Creating ${policy_name} password policy"
    cat > "${policy_name}.hcl" << EOF
length = ${length}

rule "charset" {
  charset = "abcdefghijklmnopqrstuvwxyz"
  min-chars = ${lowercase_min}
}

rule "charset" {
  charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  min-chars = ${uppercase_min}
}

rule "charset" {
  charset = "0123456789"
  min-chars = ${digit_min}
}

rule "charset" {
  charset = "${special_chars}"
  min-chars = ${special_min}
}
EOF

    "$binpath" write sys/policies/password/"${policy_name}" policy=@"${policy_name}.hcl"
    echo "Password policy '${policy_name}' created"
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
[[ -z "$STRONG_POLICY" ]] && fail "STRONG_POLICY env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

CREDENTIAL_TTL_BUFFER=${CREDENTIAL_TTL_BUFFER:-80}
DYNAMIC_ROLE_NAME="dynamic-role"
INVALID_ROLE_NAME="invalid-role"

# Define password policy constants
readonly DEFAULT_POLICY="default-policy"
readonly DEFAULT_LENGTH=20
readonly DEFAULT_MIN_CHARS=1
readonly STRONG_POLICY_LENGTH=24
readonly SPECIAL_CHARS="!@#\$%^&*"

# Create password policy for LDAP
create_password_policy "$DEFAULT_POLICY" "$DEFAULT_LENGTH" "$DEFAULT_MIN_CHARS" "$DEFAULT_MIN_CHARS" "$DEFAULT_MIN_CHARS" "$SPECIAL_CHARS" "$DEFAULT_MIN_CHARS"

# Set LDAP bind DN, password, and server URL for Vault's LDAP secrets engine
echo "Vault: Configuring LDAP secrets engine"
"$binpath" write "${MOUNT}/config" \
  binddn="cn=admin,dc=${LDAP_USERNAME},dc=com" \
  bindpass="${LDAP_ADMIN_PW}" \
  url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
  userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
  schema="openldap" \
  password_policy="$DEFAULT_POLICY"

echo "LDAP config written with $DEFAULT_POLICY"

# Create LDIF files for dynamic LDAP role
echo "Vault: Creating creation.ldif for dynamic LDAP role"
CREATION_LDIF="creation.ldif"
generate_ldif "${CREATION_LDIF}" "objectClass: person
objectClass: top
cn: {{.Username}}
sn: {{.Password | utf16le | base64}}
userPassword: {{.Password}}" "users"

echo "Vault: Creating deletion.ldif for dynamic LDAP role"
DELETION_LDIF="deletion.ldif"
generate_ldif "${DELETION_LDIF}" "changetype: delete" "users"

echo "Vault: Creating rollback.ldif for dynamic LDAP role"
ROLLBACK_LDIF="rollback.ldif"
generate_ldif "${ROLLBACK_LDIF}" "changetype: delete" "users"

# Create the dynamic LDAP role in Vault using the LDIF templates
echo "Vault: Creating dynamic LDAP role ${DYNAMIC_ROLE_NAME}"
"$binpath" write "${MOUNT}/role/${DYNAMIC_ROLE_NAME}" \
  creation_ldif=@${CREATION_LDIF} \
  deletion_ldif=@${DELETION_LDIF} \
  rollback_ldif=@${ROLLBACK_LDIF} \
  default_ttl="${DEFAULT_TTL}" \
  max_ttl="${MAX_TTL}"

# Test1: Read Root Credential Config - Verify password is excluded
test_read_config_password_exclusion() {
  echo "Test: Reading LDAP config to verify password exclusion"
  CONFIG_OUTPUT=$("$binpath" read "${MOUNT}/config")

  BIND_DN=$(jq -r '.data.binddn' <<< "$CONFIG_OUTPUT")
  URL=$(jq -r '.data.url' <<< "$CONFIG_OUTPUT")
  SCHEMA=$(jq -r '.data.schema' <<< "$CONFIG_OUTPUT")
  PASSWORD_POLICY=$(jq -r '.data.password_policy' <<< "$CONFIG_OUTPUT")

  # Verify bindpass is NOT returned (should not exist in the response)
  if ! jq -e '.data.bindpass' <<< "$CONFIG_OUTPUT" > /dev/null 2>&1; then
    echo "SUCCESS: bindpass is excluded from config read"
  else
    fail "ERROR: bindpass was returned in config"
  fi

  # Verify other expected fields are present
  if [[ -n "$BIND_DN" ]] && [[ -n "$URL" ]] && [[ -n "$SCHEMA" ]]; then
    echo "SUCCESS: All config fields are present when reading config"
  else
    fail "ERROR: Some config fields are missing when reading config"
  fi

  # Verify password policy is set
  if [[ "$PASSWORD_POLICY" == "$DEFAULT_POLICY" ]]; then
    echo "SUCCESS: Password policy '$DEFAULT_POLICY' is configured"
  else
    fail "ERROR: Password policy not set correctly. Expected: $DEFAULT_POLICY, Got: $PASSWORD_POLICY"
  fi
}

# Test: Update Root Credentials Config - Update with different password policy
test_update_password_policy() {
  echo "Test: Updating LDAP config with different password policy"

  # Create a different password policy (strong-policy)
  create_password_policy "$STRONG_POLICY" "$STRONG_POLICY_LENGTH" "$DEFAULT_MIN_CHARS" "$DEFAULT_MIN_CHARS" "$DEFAULT_MIN_CHARS" "$SPECIAL_CHARS" "$DEFAULT_MIN_CHARS"

  "$binpath" write "${MOUNT}/config" \
    password_policy="$STRONG_POLICY"

  echo "Updated LDAP config with $STRONG_POLICY"

  # Verify the update by reading config again
  UPDATED_CONFIG=$("$binpath" read "${MOUNT}/config")
  UPDATED_POLICY=$(jq -r '.data.password_policy' <<< "$UPDATED_CONFIG")

  if [[ "$UPDATED_POLICY" == "$STRONG_POLICY" ]]; then
    echo "SUCCESS: Password policy updated to '$STRONG_POLICY'"
  else
    fail "ERROR: Password policy not updated. Expected: $STRONG_POLICY, Got: $UPDATED_POLICY"
  fi
}

# Test: Generate and verify dynamic LDAP credentials
test_dynamic_credentials() {

  # Generating dynamic LDAP credentials from role ${DYNAMIC_ROLE_NAME}
  echo "Vault: Generating dynamic LDAP credentials from role ${DYNAMIC_ROLE_NAME}"

  DYNAMIC_CREDS=$("$binpath" read "${MOUNT}/creds/${DYNAMIC_ROLE_NAME}")

  DYNAMIC_PASSWORD=$(jq -r '.data.password' <<< "$DYNAMIC_CREDS")
  DYN_DN=$(jq -r '.data.distinguished_names[0]' <<< "$DYNAMIC_CREDS")

  if [[ -z "$DYN_DN" || "$DYN_DN" == "null" ]]; then
    fail "Vault did not return a distinguished_name for dynamic LDAP user"
  fi

  # Verify Dynamic Credentials
  if ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "${DYN_DN}" -w "${DYNAMIC_PASSWORD}" > /dev/null 2>&1; then
    echo "LDAP dynamic credentials valid"
  else
    fail "Error: LDAP dynamic credentials validation failed"
  fi
}

# Test: Verify credentials expire after TTL
test_credential_expiration() {

  sleep "$CREDENTIAL_TTL_BUFFER"

  if ldapwhoami -x -H "ldap://${LDAP_SERVER}:${LDAP_PORT}" -D "$DYN_DN" -w "$DYNAMIC_PASSWORD" &> /dev/null; then
    fail "Error: Dynamic credentials still valid — TTL did NOT expire"
  else
    echo "Dynamic credentials expired as expected"
  fi
}

# Test: Invalid dynamic LDAP role configuration
test_invalid_role_configuration() {
  echo "Test: Testing invalid dynamic LDAP role configuration"

  INVALID_CREATION_LDIF="invalid_creation.ldif"
  INVALID_DELETION_LDIF="invalid_deletion.ldif"

  generate_ldif "${INVALID_CREATION_LDIF}" "objectClass: person
objectClass: top
cn: {{.Username}}
sn: {{.Password}}
userPassword: {{.Password}}" "invalid_ou"

  generate_ldif "${INVALID_DELETION_LDIF}" "changetype: delete" "invalid_ou"

  echo "Vault: Attempting to create ${INVALID_ROLE_NAME} (should fail)..."

  "$binpath" write "${MOUNT}/role/${INVALID_ROLE_NAME}" \
    creation_ldif=@${INVALID_CREATION_LDIF} \
    deletion_ldif=@${INVALID_DELETION_LDIF} \
    rollback_ldif=@${INVALID_DELETION_LDIF} \
    default_ttl="${DEFAULT_TTL}" \
    max_ttl="${MAX_TTL}" || true

  echo "Attempting to read creds from ${INVALID_ROLE_NAME}"

  if ! "$binpath" read "${MOUNT}/creds/${INVALID_ROLE_NAME}" > /dev/null 2>&1; then
    echo "SUCCESS: Vault failed dynamic credential creation due to invalid OU/DN."
    "$binpath" delete "${MOUNT}/role/${INVALID_ROLE_NAME}"
  else
    fail "ERROR: Vault did NOT fail when invalid DN/OU was used!"
  fi
}

# Run test cases
test_read_config_password_exclusion
test_update_password_policy
test_dynamic_credentials
test_credential_expiration
test_invalid_role_configuration

echo "All tests completed successfully!"
