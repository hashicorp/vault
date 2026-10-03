#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

normalize_ttl() {
  case "$1" in
    null | "") echo 0 ;;
    *h) echo $((${1%h} * 3600))   ;;
    *m) echo $((${1%m} * 60))   ;;
    *s) echo $((${1%s}))   ;;
    *) echo "$1" ;; # assume already in seconds
  esac
}

log() {
  echo "[DEBUG] $1" >&2
}

log "Starting env var checks"

# Common required vars
[[ -z "$VERIFY_SSH_SECRETS" ]] && fail "VERIFY_SSH_SECRETS env variable has not been set"
[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$KEY_TYPE" ]] && fail "KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

# Always required for both types
[[ -z "$DEFAULT_USER" ]] && fail "DEFAULT_USER env variable has not been set"
[[ -z "$ALLOWED_USERS" ]] && fail "ALLOWED_USERS env variable has not been set"

if [[ "$VERIFY_SSH_SECRETS" == "false" ]]; then
  log "VERIFY_SSH_SECRETS is false; exiting script"
  exit 0
fi

# Type-specific required vars
case "$KEY_TYPE" in
  otp)
    [[ -z "$PORT" ]] && fail "PORT env variable has not been set"
    [[ -z "$CIDR_LIST" ]] && fail "CIDR_LIST env variable has not been set"
    [[ -z "$EXCLUDE_CIDR_LIST" ]] && fail "EXCLUDE_CIDR_LIST env variable has not been set"
    ;;
  ca)
    [[ -z "$TTL" ]] && fail "TTL env variable has not been set"
    [[ -z "$MAX_TTL" ]] && fail "MAX_TTL env variable has not been set"
    [[ -z "$KEY_ID_FORMAT" ]] && fail "KEY_ID_FORMAT env variable has not been set"
    [[ -z "$ALLOW_USER_CERTIFICATES" ]] && fail "ALLOW_USER_CERTIFICATES env variable has not been set"
    [[ -z "$ALLOW_HOST_CERTIFICATES" ]] && fail "ALLOW_HOST_CERTIFICATES env variable has not been set"
    [[ -z "$ALLOW_USER_KEY_IDS" ]] && fail "ALLOW_USER_KEY_IDS env variable has not been set"
    [[ -z "$ALLOW_EMPTY_PRINCIPALS" ]] && fail "ALLOW_EMPTY_PRINCIPALS env variable has not been set"
    [[ -z "$ALGORITHM_SIGNER" ]] && fail "ALGORITHM_SIGNER env variable has not been set"
    ;;
  *)
    fail "Unsupported KEY_TYPE in env check: $KEY_TYPE"
    ;;
esac

log "finished env var checks"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" read "ssh/roles/$ROLE_NAME" 2>&1); then
  fail "failed to read ssh/roles/$ROLE_NAME: $output"
fi

log "Successfully read role $ROLE_NAME"

key_type=$(jq -r '.data.key_type' <<< "$output")
default_user=$(jq -r '.data.default_user' <<< "$output")
allowed_users=$(jq -r '.data.allowed_users' <<< "$output")

log "extracted common data"

case "$KEY_TYPE" in
  otp)
    port=$(jq -r '.data.port' <<< "$output")
    cidr_list=$(jq -r '.data.cidr_list' <<< "$output")
    exclude_cidr_list=$(jq -r '.data.exclude_cidr_list' <<< "$output")
    log "extracted otp specific data"
    ;;
  ca)
    ttl=$(jq -r '.data.ttl' <<< "$output")
    max_ttl=$(jq -r '.data.max_ttl' <<< "$output")
    key_id_format=$(jq -r '.data.key_id_format' <<< "$output")
    allow_user_certificates=$(jq -r '.data.allow_user_certificates' <<< "$output")
    allow_host_certificates=$(jq -r '.data.allow_host_certificates' <<< "$output")
    allow_user_key_ids=$(jq -r '.data.allow_user_key_ids' <<< "$output")
    allow_empty_principals=$(jq -r '.data.allow_empty_principals' <<< "$output")
    algorithm_signer=$(jq -r '.data.algorithm_signer' <<< "$output")
    log "extracted ca specific data"
    ;;
  *)
    fail "Unsupported KEY_TYPE: $KEY_TYPE"
    ;;
esac

# Verify
[[ "$key_type" != "$KEY_TYPE" ]] && fail "Key type mismatch: expected $KEY_TYPE, got $key_type"
[[ "$default_user" != "$DEFAULT_USER" ]] && fail "Default user mismatch: expected $DEFAULT_USER, got $default_user"
[[ "$allowed_users" != "$ALLOWED_USERS" ]] && fail "Allowed users mismatch: expected $ALLOWED_USERS, got $allowed_users"

log "verified common data"

case "$KEY_TYPE" in
  otp)
    [[ "$port" != "$PORT" ]] && fail "Port mismatch: expected $PORT, got $port"
    [[ "$cidr_list" != "$CIDR_LIST" ]] && fail "CIDR list mismatch: expected $CIDR_LIST, got $cidr_list"
    [[ "$exclude_cidr_list" != "$EXCLUDE_CIDR_LIST" ]] && fail "Exclude CIDR list mismatch: expected $EXCLUDE_CIDR_LIST, got $exclude_cidr_list"
    log "verified otp specific data"
    ;;
  ca)
    [[ "$(normalize_ttl "$ttl")" != "$(normalize_ttl "$TTL")" ]] && fail "TTL mismatch: expected $TTL, got $ttl"
    [[ "$(normalize_ttl "$max_ttl")" != "$(normalize_ttl "$MAX_TTL")" ]] && fail "Max TTL mismatch: expected $MAX_TTL, got $max_ttl"
    [[ "$key_id_format" != "$KEY_ID_FORMAT" ]] && fail "Key ID format mismatch: expected $KEY_ID_FORMAT, got $key_id_format"
    [[ "$allow_user_certificates" != "$ALLOW_USER_CERTIFICATES" ]] && fail "Allow user certificates mismatch: expected $ALLOW_USER_CERTIFICATES, got $allow_user_certificates"
    [[ "$allow_host_certificates" != "$ALLOW_HOST_CERTIFICATES" ]] && fail "Allow host certificates mismatch: expected $ALLOW_HOST_CERTIFICATES, got $allow_host_certificates"
    [[ "$allow_user_key_ids" != "$ALLOW_USER_KEY_IDS" ]] && fail "Allow user key IDs mismatch: expected $ALLOW_USER_KEY_IDS, got $allow_user_key_ids"
    [[ "$allow_empty_principals" != "$ALLOW_EMPTY_PRINCIPALS" ]] && fail "Allow empty principals mismatch: expected $ALLOW_EMPTY_PRINCIPALS, got $allow_empty_principals"
    [[ "$algorithm_signer" != "$ALGORITHM_SIGNER" ]] && fail "Algorithm signer mismatch: expected $ALGORITHM_SIGNER, got $algorithm_signer"
    log "verified ca specific data"
    ;;
  *)
    fail "Unsupported KEY_TYPE in verification: $KEY_TYPE"
    ;;
esac
