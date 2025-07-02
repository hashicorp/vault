#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

normalize_ttl() {
  case "$1" in
    *h) echo $((${1%h} * 3600))   ;;
    *m) echo $((${1%m} * 60))   ;;
    *s) echo $((${1%s}))   ;;
    *) echo "$1" ;; # assume already in seconds
  esac
}

# Common required vars
[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$KEY_TYPE" ]] && fail "KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

# Always required for both types
[[ -z "$DEFAULT_USER" ]] && fail "DEFAULT_USER env variable has not been set"
[[ -z "$ALLOWED_USERS" ]] && fail "ALLOWED_USERS env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$MAX_TTL" ]] && fail "MAX_TTL env variable has not been set"
[[ -z "$PORT" ]] && fail "PORT env variable has not been set"

# Type-specific required vars
if [[ "$KEY_TYPE" == "otp" ]]; then
  [[ -z "$CIDR_LIST" ]] && fail "CIDR_LIST env variable has not been set"
  [[ -z "$EXCLUDE_CIDR_LIST" ]] && fail "EXCLUDE_CIDR_LIST env variable has not been set"
elif [[ "$KEY_TYPE" == "ca" ]]; then
  [[ -z "$KEY_ID_FORMAT" ]] && fail "KEY_ID_FORMAT env variable has not been set"
  [[ -z "$ALLOWED_EXTENSIONS" ]] && fail "ALLOWED_EXTENSIONS env variable has not been set"
  [[ -z "$DEFAULT_EXTENSIONS" ]] && fail "DEFAULT_EXTENSIONS env variable has not been set"
  [[ -z "$ALLOW_USER_CERTIFICATES" ]] && fail "ALLOW_USER_CERTIFICATES env variable has not been set"
  [[ -z "$ALLOW_HOST_CERTIFICATES" ]] && fail "ALLOW_HOST_CERTIFICATES env variable has not been set"
  [[ -z "$ALLOW_USER_KEY_IDS" ]] && fail "ALLOW_USER_KEY_IDS env variable has not been set"
  [[ -z "$ALLOW_EMPTY_PRINCIPALS" ]] && fail "ALLOW_EMPTY_PRINCIPALS env variable has not been set"
  [[ -z "$ALGORITHM_SIGNER" ]] && fail "ALGORITHM_SIGNER env variable has not been set"
fi

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" read "ssh/roles/$ROLE_NAME" 2>&1); then
  fail "failed to read ssh/roles/$ROLE_NAME: $output"
fi

key_type=$(echo "$output" | jq -r '.data.key_type')
default_user=$(echo "$output" | jq -r '.data.default_user')
allowed_users=$(echo "$output" | jq -r '.data.allowed_users')
port=$(echo "$output" | jq -r '.data.port')
ttl=$(echo "$output" | jq -r '.data.ttl')
max_ttl=$(echo "$output" | jq -r '.data.max_ttl')

if [[ "$KEY_TYPE" == "otp" ]]; then
  cidr_list=$(echo "$output" | jq -r '.data.cidr_list')
  exclude_cidr_list=$(echo "$output" | jq -r '.data.exclude_cidr_list')
elif [[ "$KEY_TYPE" == "ca" ]]; then
  key_id_format=$(echo "$output" | jq -r '.data.key_id_format')
  allowed_extensions=$(echo "$output" | jq -r '.data.allowed_extensions')
  default_extensions=$(echo "$output" | jq -r '.data.default_extensions')
  allow_user_certificates=$(echo "$output" | jq -r '.data.allow_user_certificates')
  allow_host_certificates=$(echo "$output" | jq -r '.data.allow_host_certificates')
  allow_user_key_ids=$(echo "$output" | jq -r '.data.allow_user_key_ids')
  allow_empty_principals=$(echo "$output" | jq -r '.data.allow_empty_principals')
  algorithm_signer=$(echo "$output" | jq -r '.data.algorithm_signer')
fi

# Verify
[[ "$key_type" != "$KEY_TYPE" ]] && fail "Key type mismatch: expected $KEY_TYPE, got $key_type"
[[ "$default_user" != "$DEFAULT_USER" ]] && fail "Default user mismatch: expected $DEFAULT_USER, got $default_user"
[[ "$allowed_users" != "$ALLOWED_USERS" ]] && fail "Allowed users mismatch: expected $ALLOWED_USERS, got $allowed_users"
[[ "$(normalize_ttl "$ttl")" != "$(normalize_ttl "$TTL")" ]] && fail "TTL mismatch: expected $TTL, got $ttl"
[[ "$(normalize_ttl "$max_ttl")" != "$(normalize_ttl "$MAX_TTL")" ]] && fail "Max TTL mismatch: expected $MAX_TTL, got $max_ttl"
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
