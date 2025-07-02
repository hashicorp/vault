#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$SIGNED_KEY" ]] && fail "SIGNED_KEY env variable has not been set"
[[ -z "$CERT_KEY_TYPE" ]] && fail "CA_KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

SIGNED_KEY_PATH=$(mktemp)
trap 'rm -f "$SIGNED_KEY_PATH"' EXIT
echo "$SIGNED_KEY" > "$SIGNED_KEY_PATH"

# Inspect the signed key
if ! ssh_key_info=$(ssh-keygen -Lf "$SIGNED_KEY_PATH"); then
  fail "Failed to verify signed SSH key"
fi

# Extract key type
cert_key_type=$(echo "$ssh_key_info" | grep "Type:" | awk '{print $2}')
if [[ "$cert_key_type" != *"$CERT_KEY_TYPE"* ]]; then
  fail "Key type mismatch: expected $CA_KEY_TYPE, got $ca_key_type"
fi
