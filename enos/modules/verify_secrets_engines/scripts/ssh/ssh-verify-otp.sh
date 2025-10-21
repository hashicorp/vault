#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

log() {
  echo "[DEBUG] $1" >&2
}

[[ -z "$VERIFY_SSH_SECRETS" ]] && fail "VERIFY_SSH_SECRETS env variable has not been set"
[[ -z "$IP" ]] && fail "IP env variable has not been set"
[[ -z "$USERNAME" ]] && fail "USERNAME env variable has not been set"
[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

if [[ "$VERIFY_SSH_SECRETS" == "false" ]]; then
  log "VERIFY_SSH_SECRETS is false; exiting script"
  exit 0
fi

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

log "Generating OTP credential from Vault"
otp_cred=$("$binpath" write -format=json "ssh/creds/$ROLE_NAME" ip="$IP" username="$USERNAME") \
  || fail "Failed to generate OTP credential"

OTP=$(jq -r '.data.key' <<< "$otp_cred")
log "Generated OTP: $OTP"

log "Verifying OTP"
otp_output=$("$binpath" write -format=json ssh/verify otp="$OTP") \
  || fail "Failed to verify OTP credential for key $OTP"

log "OTP Verification successful"
ip=$(jq -r '.data.ip' <<< "$otp_output")
role_name=$(jq -r '.data.role_name' <<< "$otp_output")
username=$(jq -r '.data.username' <<< "$otp_output")

log "IP: $ip"
log "Role Name: $role_name"
log "Username: $username"

[[ "$ip" != "$IP" ]] && fail "IP mismatch: expected $ip, got $IP"
[[ "$role_name" != "$ROLE_NAME" ]] && fail "Role name mismatch: expected $role_name, got $ROLE_NAME"
[[ "$username" != "$USERNAME" ]] && fail "Username mismatch: expected $username, got $USERNAME"

log "Completed with no mismatches"
