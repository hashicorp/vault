#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$OTP_ROLE_NAME" ]] && fail "OTP_ROLE_NAME env variable has not been set"
[[ -z "$TARGET_IP" ]] && fail "TARGET_IP env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! otp_output=$("$binpath" write ssh/creds/"$OTP_ROLE_NAME" ip="$TARGET_IP" 2>&1); then
  fail "failed to generate OTP credential: $otp_output"
fi

otp_key=$(echo "$otp_output" | jq -r '.data.key')
otp_user=$(echo "$otp_output" | jq -r '.data.username')

echo "OTP credential generated successfully for user $otp_user with OTP: $otp_key"
