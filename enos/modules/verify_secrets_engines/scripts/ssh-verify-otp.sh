#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$IP" ]] && fail "IP env variable has not been set"
[[ -z "$USERNAME" ]] && fail "USERNAME env variable has not been set"
[[ -z "$ROLE_NAME" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

otp_cred=$("$binpath" write -format=json "ssh/creds/$ROLE_NAME" ip="$IP" username="$USERNAME") \
  || fail "Failed to generate OTP credential"

OTP=$(echo "$otp_cred" | jq -r '.data.key')

otp_output=$("$binpath" write -format=json ssh/verify otp="$OTP") \
  || fail "Failed to verify OTP credential for key $OTP"

ip=$(echo "$otp_output" | jq -r '.data.ip')
role_name=$(echo "$otp_output" | jq -r '.data.role_name')
username=$(echo "$otp_output" | jq -r '.data.username')

[[ "$ip" != "$IP" ]] && fail "IP mismatch: expected $ip, got $IP"
[[ "$role_name" != "$ROLE_NAME" ]] && fail "Role name mismatch: expected $role_name, got $ROLE_NAME"
[[ "$username" != "$USERNAME" ]] && fail "Username mismatch: expected $username, got $USERNAME"
