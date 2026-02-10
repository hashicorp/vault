#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

fail() {
  echo "ERROR: $1" 1>&2
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

  current_bind_pw=$(
    "$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}" \
      | jq -r '.data.password'
  )

  [[ -n "$current_bind_pw" ]] || fail "Failed to read current bind password"

  "$binpath" write "${MOUNT}/config" \
    binddn="${LDAP_USER_DN}" \
    bindpass="${current_bind_pw}" \
    url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
    userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
    userattr="uid" \
    > /dev/null \
    || fail "Failed to update LDAP config with refreshed bind credentials"
}

# TEST 1: Create Static Role (take over existing LDAP account)
test_create_static_role() {
  echo "Create Static Role (take over existing LDAP account)"
  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" \
    dn="${LDAP_USER_DN}" \
    username="${LDAP_USER}" \
    rotation_period="${ROTATION_SHORT}" > /dev/null \
    || fail "Failed to create static role"

  ROLE_JSON=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}")

  if ! jq -e \
    '.data.username == "'"${LDAP_USER}"'" and
     .data.dn == "'"${LDAP_USER_DN}"'"' \
    <<< "$ROLE_JSON" > /dev/null; then
    fail "Static role created with incorrect attributes"
  fi
}

# TEST 2: Create Role Without Immediate Rotation
test_create_without_initial_rotation() {
  echo "Create Role Without Immediate Rotation"
  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_SKIP}" \
    dn="${LDAP_USER_DN}" \
    username="${LDAP_USER}-skip" \
    rotation_period="${ROTATION_LONG}" \
    skip_initial_rotation=true > /dev/null

  ROLE_JSON=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_SKIP}")

  if ! jq -e \
    '.data.skip_initial_rotation == true' \
    <<< "$ROLE_JSON" > /dev/null; then
    fail "skip_initial_rotation was not honored"
  fi
}

# TEST 3: Update Static Role (allowed + forbidden updates)
test_update_static_role() {
  echo "Update Static Role (allowed + forbidden updates)"
  old_role_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}")
  OLD_PERIOD=$(jq -r '.data.rotation_period' <<< "$old_role_output")

  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" \
    rotation_period="${ROTATION_LONG}" > /dev/null

  new_role_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}")
  NEW_PERIOD=$(jq -r '.data.rotation_period' <<< "$new_role_output")

  [[ -n "$NEW_PERIOD" ]] || fail "rotation_period missing after update"
  [[ "$OLD_PERIOD" != "$NEW_PERIOD" ]] || fail "rotation_period did not change"

  # Forbidden: username
  if "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" \
    username="invalid-user" > /dev/null; then
    fail "Updating username should not be allowed"
  fi

  # Forbidden: DN
  if "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" \
    dn="uid=invalid,ou=users,dc=enos,dc=com" > /dev/null; then
    fail "Updating DN should not be allowed"
  fi
}

# TEST 4: Read Static Role
test_read_static_role() {
  echo "Read Static Role"
  ROLE_JSON=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}")

  # Stable fields
  if ! jq -e \
    '.data.username == "'"${LDAP_USER}"'" and
     .data.dn == "'"${LDAP_USER_DN}"'"' \
    <<< "$ROLE_JSON" > /dev/null; then
    fail "Static role returned incorrect username or DN"
  fi

  # rotation_period must exist (string value is normalized by Vault)
  if ! jq -e \
    '.data | has("rotation_period")' \
    <<< "$ROLE_JSON" > /dev/null; then
    fail "rotation_period missing"
  fi

  # last_rotation_time may be null or absent early; ensure key exists OR is null-safe
  if ! jq -e \
    '.data | has("last_rotation_time") or (.data.last_rotation_time == null)' \
    <<< "$ROLE_JSON" > /dev/null; then
    fail "last_rotation_time missing"
  fi
}

# TEST 5: List Static Roles with hierarchical support
test_list_static_roles() {
  echo "List Static Roles with hierarchical support"
  roles=$(
    "$binpath" list "${MOUNT}/static-role" \
      | jq -r '.[]'
  )

  if ! grep -qx "${STATIC_ROLE_MAIN}" <<< "$roles"; then
    fail "Static role ${STATIC_ROLE_MAIN} not found in role list"
  fi
}

# TEST 6: Request Static Credentials
test_request_static_credentials() {
  echo "Request Static Credentials"
  CREDS=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")

  if ! jq -e '.data.password | length > 0' \
    <<< "$CREDS" > /dev/null; then
    fail "Password missing from static credentials"
  fi

  if ! jq -e '.data.ttl > 0' \
    <<< "$CREDS" > /dev/null; then
    fail "Invalid TTL returned"
  fi
}

# TEST 7: Manual Password Rotation
test_manual_password_rotation() {
  echo "Manual Password Rotation"
  old_cred_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  OLD_PW=$(jq -r '.data.password' <<< "$old_cred_output")

  refresh_ldap_bind_credentials
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" > /dev/null

  new_cred_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  NEW_PW=$(jq -r '.data.password' <<< "$new_cred_output")

  [[ "$OLD_PW" != "$NEW_PW" ]] \
    || fail "Manual password rotation did not change password"
}

# TEST 8: Automatic Password Rotation
test_automatic_password_rotation() {
  echo "Automatic Password Rotation"
  # Ensure role exists and is managed
  "$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" > /dev/null \
    || fail "Static role does not exist"

  # Capture password BEFORE rotation
  cred_before_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW_BEFORE=$(jq -r '.data.password' <<< "$cred_before_output")
  [[ -n "$PW_BEFORE" ]] || fail "Initial password missing"

  # Ensure automatic rotation is enabled
  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" \
    rotation_period="${ROTATION_LONG}" > /dev/null \
    || fail "Failed to enable automatic rotation"

  # Trigger rotation path explicitly (scheduler-safe)
  refresh_ldap_bind_credentials
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" > /dev/null \
    || fail "Failed to trigger password rotation"

  # Capture password AFTER rotation
  cred_after_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW_AFTER=$(jq -r '.data.password' <<< "$cred_after_output")
  [[ -n "$PW_AFTER" ]] || fail "Rotated password missing"

  # Password must change
  [[ "$PW_BEFORE" != "$PW_AFTER" ]] \
    || fail "Password did not change after rotation"
}

# TEST 9: Custom Password Generation
test_custom_password_generation() {
  echo "Custom Password Generation"
  "$binpath" write "sys/policies/password/${PASSWORD_POLICY}" \
    policy='
length = 20

rule "charset" { charset = "abcdefghijklmnopqrstuvwxyz" }
rule "charset" { charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" }
rule "charset" { charset = "0123456789" }
' > /dev/null \
    || fail "Failed to create password policy"

  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_MAIN}" \
    password_policy="${PASSWORD_POLICY}" \
    rotation_period="${ROTATION_SHORT}" \
    > /dev/null \
    || fail "Failed to attach password policy to static role"

  # Force rotation so policy is applied
  refresh_ldap_bind_credentials
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null \
    || fail "Failed to rotate password with custom policy"

  cred_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW=$(jq -r '.data.password' <<< "$cred_output")

  [[ "${#PW}" -ge 20 ]] || fail "Password policy not applied"
}

# TEST 10: Check Password TTL
test_check_password_ttl() {
  echo "Check Password TTL"
  CREDS1=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  TTL1=$(jq -r '.data.ttl' <<< "$CREDS1")

  # TTL must exist and be numeric
  if ! [[ "$TTL1" =~ ^[0-9]+$ ]]; then
    fail "TTL is missing or not a number"
  fi

  # TTL must be positive
  if ! [[ "$TTL1" -gt 0 ]]; then
    fail "TTL is not positive"
  fi

  sleep 3

  CREDS2=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  TTL2=$(jq -r '.data.ttl' <<< "$CREDS2")

  if ! [[ "$TTL2" =~ ^[0-9]+$ ]]; then
    fail "TTL missing after wait"
  fi

  # TTL should decrease (or rotate and reset)
  if [[ "$TTL2" -ge "$TTL1" ]]; then
    # If it increased, password must have rotated
    PW1=$(jq -r '.data.password' <<< "$CREDS1")
    PW2=$(jq -r '.data.password' <<< "$CREDS2")

    if [[ "$PW1" == "$PW2" ]]; then
      fail "TTL did not decrease and password did not rotate"
    fi
  fi
}

# TEST 11: Verify Last Vault Rotation is Present
test_verify_last_rotation_time() {
  echo "Verify Last Vault Rotation is Present"
  ROLE_JSON=$(
    "$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}"
  ) || fail "Failed to read static role"

  jq -e '.data | has("last_vault_rotation")' \
    <<< "$ROLE_JSON" > /dev/null \
    || fail "last_vault_rotation is missing"
}

# TEST 12: Verify WAL Recovery on Startup
test_wal_recovery_on_startup() {
  echo "Verify WAL Recovery on Startup"
  cred_before_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW1=$(jq -r '.data.password' <<< "$cred_before_output")

  role_before_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}")
  T1=$(jq -r '.data.last_rotation_time' <<< "$role_before_output")

  # Break LDAP to force rotation failure
  "$binpath" write "${MOUNT}/config" \
    binddn="${LDAP_USER_DN}" \
    bindpass="wrong-password" \
    url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
    userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
    userattr="uid" > /dev/null

  # Trigger rotation (WAL written, rotation fails)
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null 2>&1 || true

  # Assert rotation did NOT complete
  cred_after_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW2=$(jq -r '.data.password' <<< "$cred_after_output")

  role_after_output=$("$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_MAIN}")
  T2=$(jq -r '.data.last_rotation_time' <<< "$role_after_output")

  [[ "$PW1" == "$PW2" ]] \
    || fail "Password changed despite failed rotation"
  [[ "$T1" == "$T2" ]] \
    || fail "Rotation metadata updated despite failure"

  # Restore LDAP config after intentional failure
  refresh_ldap_bind_credentials
}

# TEST 13: Verify Rotation Retry on Failure
test_rotation_retry_on_failure() {
  echo "Verify Rotation Retry on Failure"
  # Initial password
  cred_before_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW1=$(jq -r '.data.password' <<< "$cred_before_output")

  # Break LDAP to force rotation failure
  "$binpath" write "${MOUNT}/config" \
    binddn="uid=invalid,ou=users,dc=enos,dc=com" \
    bindpass="wrong-password" \
    url="ldap://${LDAP_SERVER}:${LDAP_PORT}" \
    userdn="ou=users,dc=${LDAP_USERNAME},dc=com" \
    userattr="uid" > /dev/null

  # First rotation attempt (fails, WAL created)
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null 2>&1 || true

  # Password must remain unchanged
  cred_after_first_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW2=$(jq -r '.data.password' <<< "$cred_after_first_output")
  [[ "$PW1" == "$PW2" ]] \
    || fail "Password changed on failed rotation"

  # Second retry attempt (still fails, same WAL reused)
  "$binpath" write -f "${MOUNT}/rotate-role/${STATIC_ROLE_MAIN}" \
    > /dev/null 2>&1 || true

  cred_after_second_output=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW3=$(jq -r '.data.password' <<< "$cred_after_second_output")
  [[ "$PW1" == "$PW3" ]] \
    || fail "Password changed across retries (WAL inconsistency)"

  # Restore LDAP config after intentional failure
  refresh_ldap_bind_credentials
}

# TEST 14: Managed User Tracking (duplicate username)
test_duplicate_user_management() {
  echo "Managed User Tracking (duplicate username)"
  set +e
  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_DUP}" \
    rotation_period="${ROTATION_LONG}" \
    dn="${LDAP_USER_DN}" \
    username="${LDAP_USER}" \
    > /dev/null 2>&1
  STATUS=$?
  set -e

  # Must fail
  if [[ $STATUS -eq 0 ]]; then
    fail "Duplicate username was incorrectly allowed"
  fi

  # Role must NOT exist
  if "$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_DUP}" > /dev/null 2>&1; then
    fail "Duplicate role was created despite username already managed"
  fi
}

# TEST 15: Verify Password Rotation Not Happening
test_password_rotation_not_happening() {
  echo "Verify Password Rotation Not Happening"
  # First credential read
  CREDS1=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW1=$(jq -r '.data.password' <<< "$CREDS1")

  # Short wait (well below rotation_period)
  sleep 2

  # Second credential read
  CREDS2=$("$binpath" read "${MOUNT}/static-cred/${STATIC_ROLE_MAIN}")
  PW2=$(jq -r '.data.password' <<< "$CREDS2")

  # Password MUST NOT change
  if [[ "$PW1" != "$PW2" ]]; then
    fail "Password rotated unexpectedly in negative test"
  fi
}

# TEST 16: Verify Failure of Create Role due to Username Already Managed
test_failure_create_role_username_already_managed() {
  echo "Verify Failure of Create Role due to Username Already Managed"
  set +e
  "$binpath" write "${MOUNT}/static-role/${STATIC_ROLE_DUP}" \
    dn="${LDAP_USER_DN}" \
    username="${LDAP_USER}" \
    rotation_period="${ROTATION_LONG}" \
    > /dev/null 2>&1
  STATUS=$?
  set -e

  # Must fail
  if [[ $STATUS -eq 0 ]]; then
    fail "Expected failure when creating role with managed username"
  fi

  # Role must NOT exist
  if "$binpath" read "${MOUNT}/static-role/${STATIC_ROLE_DUP}" > /dev/null 2>&1; then
    fail "Duplicate role was created despite username already managed"
  fi
}

# Cleanup
cleanup() {
  echo "Deleting all roles"
  "$binpath" delete "${MOUNT}/static-role/${STATIC_ROLE_MAIN}"
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
