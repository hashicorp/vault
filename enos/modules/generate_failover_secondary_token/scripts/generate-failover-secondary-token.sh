#!/usr/bin/env bash
## Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

[[ -z "${VAULT_INSTALL_DIR}" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${VAULT_TOKEN}" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "${SECONDARY_PUBLIC_KEY}" ]] && fail "SECONDARY_PUBLIC_KEY env variable has not been set"

fail() {
  echo "$1" 1>&2
  exit 1
}

binpath="${VAULT_INSTALL_DIR}"/vault
test -x "${binpath}" || fail "unable to locate vault binary at ${binpath}"

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "${end_time}" ]; do
  if secondary_token=$(${binpath} write -field token sys/replication/dr/primary/secondary-token id="${VAULT_TOKEN}" secondary_public_key="${SECONDARY_PUBLIC_KEY}"); then
    echo "${secondary_token}"
    exit 0
  fi

  sleep "${RETRY_INTERVAL}"
done

fail "Timed out trying to generate secondary token"
