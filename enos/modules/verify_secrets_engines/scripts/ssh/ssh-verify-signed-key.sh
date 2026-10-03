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
[[ -z "$SIGNED_KEY" ]] && fail "SIGNED_KEY env variable has not been set"
[[ -z "$KEY_TYPE" ]] && fail "KEY_TYPE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

if [[ "$VERIFY_SSH_SECRETS" == "false" ]]; then
  log "VERIFY_SSH_SECRETS is false; exiting script"
  exit 0
fi

SIGNED_KEY_PATH=$(mktemp)
trap 'rm -f "$SIGNED_KEY_PATH"' EXIT
echo "$SIGNED_KEY" > "$SIGNED_KEY_PATH"

# Inspect the signed key
if ! ssh_key_info=$(ssh-keygen -Lf "$SIGNED_KEY_PATH"); then
  fail "Failed to verify signed SSH key"
fi

# Extract key type
key_type=$(echo "$ssh_key_info" | grep "Type:" | awk '{print $2}')
if [[ "$key_type" != *"$KEY_TYPE"* ]]; then
  fail "Key type mismatch: expected $KEY_TYPE, got $key_type"
fi
