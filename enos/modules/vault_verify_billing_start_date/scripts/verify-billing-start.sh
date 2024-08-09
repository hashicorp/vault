#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep 30
    else
      return "$exit"
    fi
  done

  return 0
}

fail() {
  echo "$1" 1>&2
  exit 1
}

export VAULT_ADDR=http://localhost:8200
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

function enable_debugging() {
        echo "Turning debugging on.."
        export PS4='+(${BASH_SOURCE}:${LINENO})> ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
        set -x
}

enable_debugging

verify_billing_start_date() {
  # get the version of vault
  version=$("$binpath" status -format=json | jq .version)

  # Get the billing start date 
  billing_start_time=$(retry 5 "$binpath" read -format=json sys/internal/counters/config  | jq -r ".data.billing_start_timestamp")

  # Verify if the billing start date is in the latest billing year

  # macOS
  if date -v -1y > /dev/null 2>&1; then
    oneYearAgoUnix=$(TZ=UTC date -v -1y +'%s')
    billingStartUnix=$(TZ=UTC date -j -f "%Y-%m-%dT%H:%M:%SZ" "${billing_start_time}" +'%s' )
  else
  # linux and unix systems
    timeago='1 year ago'
    billingStartUnix=$(TZ=UTC date -d "$billing_start_time" +'%s')    # For "now", use $(date +'%s')
    oneYearAgoUnix=$(TZ=UTC date -d "$timeago" +'%s')
  fi

  version=$("$binpath" status -format=json | jq .version)
  vault_ps=$(pgrep vault | xargs)
  #fail "Vault ADDR: $VAULT_ADDR, Vault version: $version, Vault process: $vault_ps, Billing start date: $billing_start_time"

  if [ "$billingStartUnix" -gt "$oneYearAgoUnix" ]; then
      echo "Billing start date $billing_start_time has successfully rolled over to current year."
      exit 0
  else
        fail "On version $version, pid $vault_ps, addr $VAULT_ADDR, Billing start date $billing_start_time did not roll over to current year"
  fi
}

retry 10 verify_billing_start_date
