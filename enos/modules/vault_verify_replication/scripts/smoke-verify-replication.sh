#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_EDITION" ]] && fail "VAULT_EDITION env variable has not been set"

# Replication status endpoint should have data.mode disabled for CE release
status=$(curl "${VAULT_ADDR}/v1/sys/replication/status")
if [ "$VAULT_EDITION" == "ce" ]; then
  if [ "$(jq -r '.data.mode' <<< "$status")" != "disabled" ]; then
    fail "replication data mode is not disabled for CE release!"
  fi
else
  if [ "$(jq -r '.data.dr' <<< "$status")" == "" ]; then
    fail "DR replication should be available for an ENT release!"
  fi
  if [ "$(jq -r '.data.performance' <<< "$status")" == "" ]; then
    fail "Performance replication should be available for an ENT release!"
  fi
fi
