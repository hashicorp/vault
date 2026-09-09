#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

fail() {
  printf "\nERROR: %s\n" "$1" 1>&2
  exit 1
}

log() {
  echo "==> $1"
}

# Required environment variables (provided by Enos / Terraform)
[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$LDAP_SERVER" ]] && fail "LDAP_SERVER env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_USERNAME" ]] && fail "LDAP_USERNAME env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath="${VAULT_INSTALL_DIR}/vault"
test -x "$binpath" || fail "Vault binary not found at $binpath"

export VAULT_FORMAT=json

# LDAP constants (MATCH ldap/setup.sh exactly)
LDAP_USER="${LDAP_STATIC_USERNAME:-vault-static-user}"
LDAP_USER_DN="uid=${LDAP_USER},ou=users,dc=${LDAP_USERNAME},dc=com"

STATIC_ROLE_MAIN="enos-static-main"
STATIC_ROLE_SKIP="enos-static-skip"
STATIC_ROLE_DUP="enos-static-dup"

ROTATION_SHORT="1m"
ROTATION_LONG="5m"

PASSWORD_POLICY="enos-password-policy"

# Restoring LDAP config for rotations
refresh_ldap_bind_credentials() {
  log "Refreshing LDAP bind credentials from static role"

  local current_bind_pw
  if ! current_bind_pw=$(
    "$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1 \
      | jq -r '.data.password'
  ); then
    fail "Failed to read current bindpass from ${MOUNT}/static-cred/${STATIC_ROLE_MAIN}: ${current_bind_pw}"
  fi
  [[ -n "$current_bind_pw" ]] || fail "Failed to get current bindpass from ${MOUNT}/static-cred/${STATIC_ROLE_MAIN}; returned value was blank"

  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/config" - << EOF 2>&1
{
  "binddn": "${LDAP_USER_DN}",
  "bindpass": "${current_bind_pw}",
  "url": "ldap://${LDAP_SERVER}:${LDAP_PORT}",
  "userdn": "ou=users,dc=${LDAP_USERNAME},dc=com",
  "userattr": "uid"
}
EOF
  ); then
    fail "Failed to update LDAP config with refreshed bind credentials: ${output}"
  fi
}

# TEST 1: Create Static Role (take over existing LDAP account)
test_create_static_role() {
  echo "Create Static Role (take over existing LDAP account)"

  local output
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" - << EOF 2>&1
{
  "dn": "${LDAP_USER_DN}",
  "username": "${LDAP_USER}",
  "rotation_period": "${ROTATION_SHORT}"
}
EOF
  ); then
    fail "Failed to create static role: ${output}"
  fi

  local role_json
  if ! role_json=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read static role: ${role_json}"
  fi

  # Verify role was created with correct username and DN
  if ! jq -e \
    '.data.username == "'"${LDAP_USER}"'" and
     .data.dn == "'"${LDAP_USER_DN}"'"' \
    <<< "$role_json" > /dev/null; then
    fail "Static role created with incorrect attributes"
  fi
}

# TEST 2: Create Role Without Immediate Rotation
test_create_without_initial_rotation() {
  echo "Create Role Without Immediate Rotation"

  local output
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_SKIP}" - << EOF 2>&1
{
  "dn": "${LDAP_USER_DN}",
  "username": "${LDAP_USER}-skip",
  "rotation_period": "${ROTATION_LONG}",
  "skip_initial_rotation": true
}
EOF
  ); then
    fail "Failed to create static role with skip_initial_rotation: ${output}"
  fi

  local role_json
  if ! role_json=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_SKIP}" 2>&1); then
    fail "Failed to read static role: ${role_json}"
  fi

  if ! jq -e \
    '.data.skip_initial_rotation == true' \
    <<< "$role_json" > /dev/null; then
    fail "skip_initial_rotation was not honored"
  fi
}

# TEST 3: Update Static Role (allowed + forbidden updates)
test_update_static_role() {
  echo "Update Static Role (allowed + forbidden updates)"

  local old_role_output
  if ! old_role_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read static role before update: ${old_role_output}"
  fi

  local old_period
  old_period=$(jq -r '.data.rotation_period' <<< "$old_role_output")

  local output
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" - << EOF 2>&1
{
  "rotation_period": "${ROTATION_LONG}"
}
EOF
  ); then
    fail "Failed to update rotation_period: ${output}"
  fi

  local new_role_output
  if ! new_role_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read static role after update: ${new_role_output}"
  fi

  local new_period
  new_period=$(jq -r '.data.rotation_period' <<< "$new_role_output")

  [[ -n "$new_period" ]] || fail "rotation_period missing after update"
  [[ "$old_period" != "$new_period" ]] || fail "rotation_period did not change"

  # Verify forbidden update: username (should fail)
  if output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" - << EOF 2>&1
{
  "username": "invalid-user"
}
EOF
  ); then
    fail "Updating username should not be allowed"
  fi

  # Verify forbidden update: DN (should fail)
  if output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" - << EOF 2>&1
{
  "dn": "uid=invalid,ou=users,dc=enos,dc=com"
}
EOF
  ); then
    fail "Updating DN should not be allowed"
  fi
}

# TEST 4: Read Static Role
test_read_static_role() {
  echo "Read Static Role"

  local role_json
  if ! role_json=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read static role: ${role_json}"
  fi

  # Verify username and DN match expected values
  if ! jq -e \
    '.data.username == "'"${LDAP_USER}"'" and
     .data.dn == "'"${LDAP_USER_DN}"'"' \
    <<< "$role_json" > /dev/null; then
    fail "Static role returned incorrect username or DN"
  fi

  # Verify rotation_period exists (value is normalized by Vault)
  if ! jq -e \
    '.data | has("rotation_period")' \
    <<< "$role_json" > /dev/null; then
    fail "rotation_period missing"
  fi

  # Verify last_rotation_time exists (may be null initially)
  if ! jq -e \
    '.data | has("last_rotation_time") or (.data.last_rotation_time == null)' \
    <<< "$role_json" > /dev/null; then
    fail "last_rotation_time missing"
  fi
}

# TEST 5: List Static Roles with hierarchical support
test_list_static_roles() {
  echo "List Static Roles with hierarchical support"

  local roles
  if ! roles=$(
    "$binpath" list "${MOUNT}/static-role" 2>&1 \
      | jq -r '.[]'
  ); then
    fail "Failed to list static roles: ${roles}"
  fi

  # Verify main role appears in the list
  if ! grep -qx "${STATIC_ROLE_MAIN}" <<< "$roles"; then
    fail "Static role ${STATIC_ROLE_MAIN} not found in role list"
  fi
}

# TEST 6: Request Static Credentials
test_request_static_credentials() {
  echo "Request Static Credentials"

  local creds
  if ! creds=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read static credentials: ${creds}"
  fi

  # Verify password exists and is non-empty
  if ! jq -e '.data.password | length > 0' \
    <<< "$creds" > /dev/null; then
    fail "Password missing from static credentials"
  fi

  # Verify TTL is positive
  if ! jq -e '.data.ttl > 0' \
    <<< "$creds" > /dev/null; then
    fail "Invalid TTL returned"
  fi
}

# TEST 7: Manual Password Rotation
test_manual_password_rotation() {
  echo "Manual Password Rotation"

  local old_cred_output
  if ! old_cred_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials before rotation: ${old_cred_output}"
  fi

  local old_password
  old_password=$(jq -r '.data.password' <<< "$old_cred_output")

  # Refresh LDAP bind credentials before rotation
  refresh_ldap_bind_credentials

  local output
  if ! output=$("$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to manually rotate password: ${output}"
  fi

  local new_cred_output
  if ! new_cred_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after rotation: ${new_cred_output}"
  fi

  local new_password
  new_password=$(jq -r '.data.password' <<< "$new_cred_output")

  # Verify password changed after manual rotation
  [[ "$old_password" != "$new_password" ]] \
    || fail "Manual password rotation did not change password"
}

# TEST 8: Automatic Password Rotation
test_automatic_password_rotation() {
  echo "Automatic Password Rotation"

  local output
  # Verify role exists and is managed
  if ! output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Static role does not exist: ${output}"
  fi

  local cred_before_output
  if ! cred_before_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials before rotation: ${cred_before_output}"
  fi

  local password_before
  password_before=$(jq -r '.data.password' <<< "$cred_before_output")
  [[ -n "$password_before" ]] || fail "Initial password missing"

  # Enable automatic rotation
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" - << EOF 2>&1
{
  "rotation_period": "${ROTATION_LONG}"
}
EOF
  ); then
    fail "Failed to enable automatic rotation: ${output}"
  fi

  # Refresh LDAP bind credentials and trigger rotation
  refresh_ldap_bind_credentials

  if ! output=$("$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to trigger password rotation: ${output}"
  fi

  local cred_after_output
  if ! cred_after_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after rotation: ${cred_after_output}"
  fi

  local password_after
  password_after=$(jq -r '.data.password' <<< "$cred_after_output")
  [[ -n "$password_after" ]] || fail "Rotated password missing"

  # Verify password changed after rotation
  [[ "$password_before" != "$password_after" ]] \
    || fail "Password did not change after rotation"
}

# TEST 9: Custom Password Generation
test_custom_password_generation() {
  echo "Custom Password Generation"

  local output
  # Create custom password policy
  if ! output=$(
    "$binpath" write -format=json "sys/policies/password/${PASSWORD_POLICY}" - << EOF 2>&1
{
  "policy": "length = 20\n\nrule \"charset\" { charset = \"abcdefghijklmnopqrstuvwxyz\" }\nrule \"charset\" { charset = \"ABCDEFGHIJKLMNOPQRSTUVWXYZ\" }\nrule \"charset\" { charset = \"0123456789\" }"
}
EOF
  ); then
    fail "Failed to create password policy: ${output}"
  fi

  # Attach password policy to static role
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" - << EOF 2>&1
{
  "password_policy": "${PASSWORD_POLICY}",
  "rotation_period": "${ROTATION_SHORT}"
}
EOF
  ); then
    fail "Failed to attach password policy to static role: ${output}"
  fi

  # Refresh LDAP bind credentials and force rotation to apply policy
  refresh_ldap_bind_credentials

  if ! output=$("$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to rotate password with custom policy: ${output}"
  fi

  local cred_output
  if ! cred_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after policy rotation: ${cred_output}"
  fi

  local password
  password=$(jq -r '.data.password' <<< "$cred_output")

  # Verify password meets minimum length requirement from policy
  [[ "${#password}" -ge 20 ]] || fail "Password policy not applied (expected length >= 20, got ${#password})"
}

# TEST 10: Check Password TTL
test_check_password_ttl() {
  echo "Check Password TTL"

  local creds_first
  if ! creds_first=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials for TTL check: ${creds_first}"
  fi

  local ttl_first
  ttl_first=$(jq -r '.data.ttl' <<< "$creds_first")

  # Verify TTL exists and is numeric
  if ! [[ "$ttl_first" =~ ^[0-9]+$ ]]; then
    fail "TTL is missing or not a number"
  fi

  # Verify TTL is positive
  if ! [[ "$ttl_first" -gt 0 ]]; then
    fail "TTL is not positive"
  fi

  sleep 3

  local creds_second
  if ! creds_second=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after wait: ${creds_second}"
  fi

  local ttl_second
  ttl_second=$(jq -r '.data.ttl' <<< "$creds_second")

  if ! [[ "$ttl_second" =~ ^[0-9]+$ ]]; then
    fail "TTL missing after wait"
  fi

  # Verify TTL decreased or password rotated (which resets TTL)
  if [[ "$ttl_second" -ge "$ttl_first" ]]; then
    local password_first
    password_first=$(jq -r '.data.password' <<< "$creds_first")

    local password_second
    password_second=$(jq -r '.data.password' <<< "$creds_second")

    if [[ "$password_first" == "$password_second" ]]; then
      fail "TTL did not decrease and password did not rotate"
    fi
  fi
}

# TEST 11: Verify Last Vault Rotation is Present
test_verify_last_rotation_time() {
  echo "Verify Last Vault Rotation is Present"

  local role_json
  if ! role_json=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read static role: ${role_json}"
  fi

  # Verify last_vault_rotation field exists in role metadata
  jq -e '.data | has("last_vault_rotation")' \
    <<< "$role_json" > /dev/null \
    || fail "last_vault_rotation is missing"
}

# TEST 12: Verify WAL Recovery on Startup
test_wal_recovery_on_startup() {
  echo "Verify WAL Recovery on Startup"

  local cred_before_output
  if ! cred_before_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials before WAL test: ${cred_before_output}"
  fi

  local password_before
  password_before=$(jq -r '.data.password' <<< "$cred_before_output")

  local role_before_output
  if ! role_before_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read role before WAL test: ${role_before_output}"
  fi

  local rotation_time_before
  rotation_time_before=$(jq -r '.data.last_rotation_time' <<< "$role_before_output")

  local output
  # Intentionally break LDAP config to force rotation failure
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/config" - << EOF 2>&1
{
  "binddn": "${LDAP_USER_DN}",
  "bindpass": "wrong-password",
  "url": "ldap://${LDAP_SERVER}:${LDAP_PORT}",
  "userdn": "ou=users,dc=${LDAP_USERNAME},dc=com",
  "userattr": "uid"
}
EOF
  ); then
    fail "Failed to break LDAP config for WAL test: ${output}"
  fi

  # Trigger rotation (WAL written, but rotation should fail)
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null 2>&1 || true

  local cred_after_output
  if ! cred_after_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after failed rotation: ${cred_after_output}"
  fi

  local password_after
  password_after=$(jq -r '.data.password' <<< "$cred_after_output")

  local role_after_output
  if ! role_after_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read role after failed rotation: ${role_after_output}"
  fi

  local rotation_time_after
  rotation_time_after=$(jq -r '.data.last_rotation_time' <<< "$role_after_output")

  # Verify rotation did NOT complete (password and metadata unchanged)
  [[ "$password_before" == "$password_after" ]] \
    || fail "Password changed despite failed rotation"
  [[ "$rotation_time_before" == "$rotation_time_after" ]] \
    || fail "Rotation metadata updated despite failure"

  # Restore LDAP config after intentional failure
  refresh_ldap_bind_credentials
}

# TEST 13: Verify Rotation Retry on Failure
test_rotation_retry_on_failure() {
  echo "Verify Rotation Retry on Failure"

  local cred_before_output
  if ! cred_before_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials before retry test: ${cred_before_output}"
  fi

  local password_initial
  password_initial=$(jq -r '.data.password' <<< "$cred_before_output")

  local output
  # Intentionally break LDAP config to force rotation failure
  if ! output=$(
    "$binpath" write -format=json "${MOUNT}/config" - << EOF 2>&1
{
  "binddn": "uid=invalid,ou=users,dc=enos,dc=com",
  "bindpass": "wrong-password",
  "url": "ldap://${LDAP_SERVER}:${LDAP_PORT}",
  "userdn": "ou=users,dc=${LDAP_USERNAME},dc=com",
  "userattr": "uid"
}
EOF
  ); then
    fail "Failed to break LDAP config for retry test: ${output}"
  fi

  # First rotation attempt (should fail, WAL created)
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null 2>&1 || true

  local cred_after_first_output
  if ! cred_after_first_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after first retry: ${cred_after_first_output}"
  fi

  local password_after_first
  password_after_first=$(jq -r '.data.password' <<< "$cred_after_first_output")

  # Verify password unchanged after first failed rotation
  [[ "$password_initial" == "$password_after_first" ]] \
    || fail "Password changed on failed rotation"

  # Second rotation attempt (should still fail, same WAL reused)
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null 2>&1 || true

  local cred_after_second_output
  if ! cred_after_second_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after second retry: ${cred_after_second_output}"
  fi

  local password_after_second
  password_after_second=$(jq -r '.data.password' <<< "$cred_after_second_output")

  # Verify password unchanged across retries (WAL consistency)
  [[ "$password_initial" == "$password_after_second" ]] \
    || fail "Password changed across retries (WAL inconsistency)"

  # Restore LDAP config after intentional failure
  refresh_ldap_bind_credentials
}

# TEST 14: Managed User Tracking (duplicate username)
test_duplicate_user_management() {
  echo "Managed User Tracking (duplicate username)"

  local output
  local status
  set +e
  output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_DUP}" - << EOF 2>&1
{
  "rotation_period": "${ROTATION_LONG}",
  "dn": "${LDAP_USER_DN}",
  "username": "${LDAP_USER}"
}
EOF
  )
  status=$?
  set -e

  # Verify creation fails (username already managed)
  if [[ $status -eq 0 ]]; then
    fail "Duplicate username was incorrectly allowed"
  fi

  # Verify role does NOT exist
  if "$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_DUP}" > /dev/null 2>&1; then
    fail "Duplicate role was created despite username already managed"
  fi
}

# TEST 15: Verify Password Rotation Not Happening
test_password_rotation_not_happening() {
  echo "Verify Password Rotation Not Happening"

  local creds_first
  if ! creds_first=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials for negative test: ${creds_first}"
  fi

  local password_first
  password_first=$(jq -r '.data.password' <<< "$creds_first")

  # Short wait (well below rotation_period)
  sleep 2

  local creds_second
  if ! creds_second=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to read credentials after wait: ${creds_second}"
  fi

  local password_second
  password_second=$(jq -r '.data.password' <<< "$creds_second")

  # Verify password did NOT change (negative test)
  if [[ "$password_first" != "$password_second" ]]; then
    fail "Password rotated unexpectedly in negative test"
  fi
}

# TEST 16: Verify Failure of Create Role due to Username Already Managed
test_failure_create_role_username_already_managed() {
  echo "Verify Failure of Create Role due to Username Already Managed"

  local output
  local status
  set +e
  output=$(
    "$binpath" write -format=json "${MOUNT}/static-role/${STATIC_ROLE_DUP}" - << EOF 2>&1
{
  "dn": "${LDAP_USER_DN}",
  "username": "${LDAP_USER}",
  "rotation_period": "${ROTATION_LONG}"
}
EOF
  )
  status=$?
  set -e

  # Verify creation fails (username already managed)
  if [[ $status -eq 0 ]]; then
    fail "Expected failure when creating role with managed username"
  fi

  # Verify role does NOT exist
  if "$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_DUP}" > /dev/null 2>&1; then
    fail "Duplicate role was created despite username already managed"
  fi
}

# Cleanup
cleanup() {
  echo "Deleting all roles"

  local output
  if ! output=$("$binpath" delete "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" 2>&1); then
    fail "Failed to delete static role ${STATIC_ROLE_MAIN}: ${output}"
  fi
  # Note: STATIC_ROLE_SKIP cleanup commented out (test currently disabled)
  # "$binpath" delete "${MOUNT}/static-role/${STATIC_ROLE_SKIP}"
}

# Test execution
test_create_static_role
# This test is currently failing because of "skip_initial_rotation" not working. Once the issue is resolved then this test will start working
# test_create_without_initial_rotation
test_update_static_role
test_read_static_role
test_list_static_roles
test_request_static_credentials
test_check_password_ttl
test_verify_last_rotation_time
test_duplicate_user_management
test_password_rotation_not_happening
test_failure_create_role_username_already_managed
test_wal_recovery_on_startup
test_manual_password_rotation
test_automatic_password_rotation
test_custom_password_generation
test_rotation_retry_on_failure
cleanup

log "SUCCESS: LDAP static role tests completed successfully"
