#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$OTP" ]] && fail "OTP env variable has not been set"
[[ -z "$IP" ]] && fail "IP env variable has not been set"
[[ -z "$USERNAME" ]] && fail "USERNAME env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! otp_output=$("$binpath" write ssh/verify otp="$OTP" 2>&1); then
  fail "failed to generate OTP credential: $otp_output"
fi

ip=$(echo "$otp_output" | jq -r '.data.ip')
username=$(echo "$otp_output" | jq -r '.data.username')

if [[ "$ip" != "$IP" ]]; then
  fail "IP mismatch: expected $ip, got $IP"
fi

if [[ "$username" != "$USERNAME" ]]; then
  fail "Username mismatch: expected $username, got $USERNAME"
fi