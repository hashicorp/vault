#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${REQPATH}" ]] && fail "REQPATH env variable has not been set"
[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${VAULT_INSTALL_DIR}" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "${VAULT_TOKEN}" ]] && fail "VAULT_TOKEN env variable has not been set"

# PAYLOAD is optional - if not set or empty, use -f flag instead
binpath=${VAULT_INSTALL_DIR}/vault
test -x "${binpath}" || fail "unable to locate vault binary at ${binpath}"

export VAULT_FORMAT=json

# Check if PAYLOAD is empty or unset
if [[ -z "${PAYLOAD}" ]]; then
  # Use -f flag for empty payload (force write without data)
  if output=$("${binpath}" write -f "${REQPATH}" 2>&1); then
    printf "%s\n" "${output}"
  else
    fail "failed to write with empty payload: path=${REQPATH} out=${output}"
  fi
else
  # Use original logic with payload data
  if output=$("${binpath}" write "${REQPATH}" - <<< "${PAYLOAD}" 2>&1); then
    printf "%s\n" "${output}"
  else
    fail "failed to write payload: path=${REQPATH} payload=${PAYLOAD} out=${output}"
  fi
fi
