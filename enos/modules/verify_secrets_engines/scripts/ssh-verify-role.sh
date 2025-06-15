#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$KEY_TYPE" ]] && fail "KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" read "ssh/roles/$ROLE_NAME" 2>&1); then
  fail "failed to read ssh/roles/$ROLE_NAME: $output"
fi

key_type=$(echo "$output" | jq -r '.data.key_type')
default_user=$(echo "$output" | jq -r '.data.default_user')
allowed_users=$(echo "$output" | jq -r '.data.allowed_users')
cidr_list=$(echo "$output" | jq -r '.data.cidr_list')
port=$(echo "$output" | jq -r '.data.port')
default_user_template=$(echo "$output" | jq -r '.data.default_user_template')
allowed_users_template=$(echo "$output" | jq -r '.data.allowed_users_template')
exclude_cidr_list=$(echo "$output" | jq -r '.data.exclude_cidr_list')
ttl=$(echo "$output" | jq -r '.data.ttl')
max_ttl=$(echo "$output" | jq -r '.data.max_ttl')

# Verify
[[ "$key_type" != "$KEY_TYPE" ]] && fail "Key type mismatch: expected $KEY_TYPE, got $key_type"
[[ "$default_user" != "$DEFAULT_USER" ]] && fail "Default user mismatch: expected $DEFAULT_USER, got $default_user"
[[ "$default_user_template" != "$DEFAULT_USER_TEMPLATE" ]] && fail "Default user template mismatch: expected $DEFAULT_USER_TEMPLATE, got $default_user_template"
[[ "$allowed_users" != "$ALLOWED_USERS" ]] && fail "Allowed users mismatch: expected $ALLOWED_USERS, got $allowed_users"
[[ "$allowed_users_template" != "$ALLOWED_USERS_TEMPLATE" ]] && fail "Allowed users template mismatch: expected $ALLOWED_USERS_TEMPLATE, got $allowed_users_template"
[[ "$ttl" != "$TTL" ]] && fail "TTL mismatch: expected $TTL, got $ttl"
[[ "$max_ttl" != "$MAX_TTL" ]] && fail "Max TTL mismatch: expected $MAX_TTL, got $max_ttl"
[[ "$port" != "$PORT" ]] && fail "Port mismatch: expected $PORT, got $port"

if [[ "$KEY_TYPE" == "otp" ]]; then
  [[ "$cidr_list" != "$CIDR_LIST" ]] && fail "CIDR list mismatch: expected $CIDR_LIST, got $cidr_list"
  [[ "$exclude_cidr_list" != "$EXCLUDE_CIDR_LIST" ]] && fail "Exclude CIDR list mismatch: expected $EXCLUDE_CIDR_LIST, got $exclude_cidr_list"
elif [[ "$KEY_TYPE" == "ca" ]]; then
  [[ "$key_id_format" != "$KEY_ID_FORMAT" ]] && fail "Key ID format mismatch: expected $KEY_ID_FORMAT, got $key_id_format"
  [[ "$allowed_extensions" != "$ALLOWED_EXTENSIONS" ]] && fail "Allowed extensions mismatch: expected $ALLOWED_EXTENSIONS, got $allowed_extensions"
  [[ "$(echo "$default_extensions" | jq -c '.')" != "$(echo "$DEFAULT_EXTENSIONS" | jq -c '.')" ]] && fail "Default extensions mismatch"
  [[ "$allow_user_certificates" != "$ALLOW_USER_CERTIFICATES" ]] && fail "Allow user certificates mismatch: expected $ALLOW_USER_CERTIFICATES, got $allow_user_certificates"
  [[ "$allow_host_certificates" != "$ALLOW_HOST_CERTIFICATES" ]] && fail "Allow host certificates mismatch: expected $ALLOW_HOST_CERTIFICATES, got $allow_host_certificates"
  [[ "$allow_user_key_ids" != "$ALLOW_USER_KEY_IDS" ]] && fail "Allow user key IDs mismatch: expected $ALLOW_USER_KEY_IDS, got $allow_user_key_ids"
  [[ "$allow_empty_principals" != "$ALLOW_EMPTY_PRINCIPALS" ]] && fail "Allow empty principals mismatch: expected $ALLOW_EMPTY_PRINCIPALS, got $allow_empty_principals"
  [[ "$algorithm_signer" != "$ALGORITHM_SIGNER" ]] && fail "Algorithm signer mismatch: expected $ALGORITHM_SIGNER, got $algorithm_signer"
fi
