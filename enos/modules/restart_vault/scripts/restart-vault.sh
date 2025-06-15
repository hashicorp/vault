#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

if ! out=$(sudo systemctl stop vault 2>&1); then
  fail "failed to stop vault: $out: $(sudo systemctl status vault)"
fi

if ! out=$(sudo systemctl daemon-reload 2>&1); then
  fail "failed to daemon-reload systemd: $out" 1>&2
fi

if ! out=$(sudo systemctl start vault 2>&1); then
  fail "failed to start vault: $out: $(sudo systemctl status vault)"
fi

count=0
retries=5
while :; do
  # Check the Vault seal status
  status=$($binpath status)
  code=$?

  if [ $code == 0 ] || [ $code == 2 ]; then
    # 0 is unsealed and 2 is running but sealed
    echo "$status"
    exit 0
  fi

  printf "Waiting for Vault cluster to be ready: status code: %s, status:\n%s\n" "$code" "$status" 2>&1

  wait=$((3 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
  else
      fail "Timed out waiting for Vault node to be ready after restart"
  fi
done
