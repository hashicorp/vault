#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

binpath="${VAULT_INSTALL_DIR}/vault"

fail() {
  echo "$1" >&2
  exit 1
}

# Check required environment variables
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$STORAGE_BACKEND" ]] && fail "STORAGE_BACKEND env variable has not been set"

# Define the policy content
policy_content() {
  cat << EOF
path "sys/replication/dr/secondary/promote" {
  capabilities = [ "update" ]
}

path "sys/replication/dr/secondary/update-primary" {
  capabilities = [ "update" ]
}
EOF
  if [ "$STORAGE_BACKEND" = "raft" ]; then
    cat << EOF
path "sys/storage/raft/autopilot/state" {
  capabilities = [ "update", "read" ]
}
EOF
  fi
}

# Write the policy
$binpath policy write dr-secondary-promotion - <<< "$(policy_content)"  &> /dev/null

# Configure the failover handler token role
$binpath write auth/token/roles/failover-handler \
  allowed_policies=dr-secondary-promotion \
  orphan=true \
  renewable=false \
  token_type=batch &> /dev/null

# Create a token for the failover handler role and output the token only
$binpath token create -field=token -role=failover-handler -ttl=8h
