#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
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

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

enable_debugging() {
  echo "Turning debugging on.."
  export PS4='+(${BASH_SOURCE}:${LINENO})> ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
  set -x
}

get_billing_start_date() {
  "$binpath" read -format=json sys/internal/counters/config  | jq -r ".data.billing_start_timestamp"
}

get_target_platform() {
  uname -s
}

# Given the date as ARGV 1, return 1 year as a unix date
verify_date_is_in_current_year() {
  local billing_start_unix
  local one_year_ago_unix

  # Verify if the billing start date is in the latest billing year
  case $(get_target_platform) in
    Linux)
      billing_start_unix=$(TZ=UTC date -d "$1" +'%s')    # For "now", use $(date +'%s')
      one_year_ago_unix=$(TZ=UTC date -d "1 year ago" +'%s')
      ;;
    Darwin)
      one_year_ago_unix=$(TZ=UTC date -v -1y +'%s')
      billing_start_unix=$(TZ=UTC date -j -f "%Y-%m-%dT%H:%M:%SZ" "${1}" +'%s')
      ;;
    *)
      fail "Unsupported target host operating system: $(get_target_platform)" 1>&2
      ;;
  esac

  if [ "$billing_start_unix" -gt "$one_year_ago_unix" ]; then
    echo "Billing start date $1 has successfully rolled over to current year."
    exit 0
  else
    local vault_ps
    vault_ps=$(pgrep vault | xargs)
    echo "On version $version, pid $vault_ps, addr $VAULT_ADDR, Billing start date $1 did not roll over to current year" 1>&2
  fi
}

verify_billing_start_date() {
  local billing_start
  billing_start=$(get_billing_start_date)

  if verify_date_is_in_current_year "$billing_start"; then
    return 0
  fi

  local version
  local vault_ps
  version=$("$binpath" status -format=json | jq .version)
  vault_ps=$(pgrep vault | xargs)
  echo "On version $version, pid $vault_ps, addr $VAULT_ADDR, Billing start date $billing_start did not roll over to current year" 1>&2
  return 1
}

enable_debugging

retry 10 verify_billing_start_date
