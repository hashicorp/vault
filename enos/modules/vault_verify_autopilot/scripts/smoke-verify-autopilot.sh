#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_AUTOPILOT_UPGRADE_STATUS" ]] && fail "VAULT_AUTOPILOT_UPGRADE_STATUS env variable has not been set"
[[ -z "$VAULT_AUTOPILOT_UPGRADE_VERSION" ]] && fail "VAULT_AUTOPILOT_UPGRADE_VERSION env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

count=0
retries=8
while :; do
  state=$($binpath read -format=json sys/storage/raft/autopilot/state)
  status="$(jq -r '.data.upgrade_info.status' <<< "$state")"
  target_version="$(jq -r '.data.upgrade_info.target_version' <<< "$state")"

  if [ "$status" = "$VAULT_AUTOPILOT_UPGRADE_STATUS" ] && [ "$target_version" = "$VAULT_AUTOPILOT_UPGRADE_VERSION" ]; then
    exit 0
  fi

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    echo "Expected autopilot status to be $VAULT_AUTOPILOT_UPGRADE_STATUS, got $status"
    echo "Expected autopilot target_version to be $VAULT_AUTOPILOT_UPGRADE_VERSION, got $target_version"
    sleep "$wait"
  else
    echo "$state"
    echo "Expected autopilot status to be $VAULT_AUTOPILOT_UPGRADE_STATUS, got $status"
    echo "Expected autopilot target_version to be $VAULT_AUTOPILOT_UPGRADE_VERSION, got $target_version"
    fail "Autopilot did not get into the correct status"
  fi
done
