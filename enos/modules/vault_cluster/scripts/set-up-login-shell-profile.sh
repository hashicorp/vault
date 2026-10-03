#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

# Determine the profile file we should write to. We only want to affect login shells and bash will
# only read one of these in ordered of precendence.
determineProfileFile() {
  if [ -f "$HOME/.bash_profile" ]; then
    printf "%s/.bash_profile\n" "$HOME"
    return 0
  fi

  if [ -f "$HOME/.bash_login" ]; then
    printf "%s/.bash_login\n" "$HOME"
    return 0
  fi

  printf "%s/.profile\n" "$HOME"
}

appendVaultProfileInformation() {
  tee -a "$1" <<< "export PATH=$PATH:$VAULT_INSTALL_DIR
export VAULT_ADDR=$VAULT_ADDR
export VAULT_TOKEN=$VAULT_TOKEN"
}

main() {
  local profile_file
  if ! profile_file=$(determineProfileFile); then
    fail "failed to determine login shell profile file location"
  fi

  # If vault_cluster is used more than once, eg: autopilot or replication, this module can
  # be called more than once. Short ciruit here if our profile is already set up.
  if grep VAULT_ADDR < "$profile_file"; then
    exit 0
  fi

  if ! appendVaultProfileInformation "$profile_file"; then
    fail "failed to write vault configuration to login shell profile"
  fi

  exit 0
}

main
