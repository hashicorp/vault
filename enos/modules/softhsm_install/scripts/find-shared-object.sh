#!/usr/bin/env bash
## Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if so=$(sudo find /usr -type f -name libsofthsm2.so -print -quit); then
    echo "$so"
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out trying to locate libsofthsm2.so shared object"
