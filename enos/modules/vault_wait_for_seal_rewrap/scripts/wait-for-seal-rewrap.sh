#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

getRewrapData() {
  $binpath read sys/sealwrap/rewrap -format=json | jq -eMc '.data'
}

waitForRewrap() {
  local data
  if ! data=$(getRewrapData); then
    echo "failed getting /v1/sys/sealwrap/rewrap data" 1>&2
    return 1
  fi

  if ! jq -e '.is_running == false' <<< "$data" &> /dev/null; then
    echo "rewrap is running" 1>&2
    return 1
  fi

  if ! jq -e '.entries.failed == 0' <<< "$data" &> /dev/null; then
    local entries
    entries=$(jq -Mc '.entries.failed' <<< "$data")
    echo "rewrap has $entries failed entries" 1>&2
    return 1
  fi

  if ! jq -e '.entries.processed == .entries.succeeded' <<< "$data" &> /dev/null; then
    local processed
    local succeeded
    processed=$(jq -Mc '.entries.processed' <<< "$data")
    succeeded=$(jq -Mc '.entries.succeeded' <<< "$data")
    echo "the number of processed entries ($processed) does not equal then number of succeeded ($succeeded)" 1>&2
    return 1
  fi

  return 0
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if waitForRewrap; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for seal rewrap to be completed. Data:\n\t$(getRewrapData)"
