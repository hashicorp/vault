#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


if test "$LICENSE" = "none"; then
  exit 0
fi

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))

    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

export VAULT_ADDR=http://localhost:8200
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

# Temporary hack until we can make the unseal resource handle legacy license
# setting. If we're running 1.8 and above then we shouldn't try to set a license.
ver=$(${BIN_PATH} version)
if [[ "$(echo "$ver" |awk '{print $2}' |awk -F'.' '{print $2}')" -ge 8 ]]; then
  exit 0
fi

retry 5 "${BIN_PATH}" write /sys/license text="$LICENSE"
