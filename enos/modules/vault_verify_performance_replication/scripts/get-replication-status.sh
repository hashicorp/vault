#!/usr/bin/env bash

set -e

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  return 1
}

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((10 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

test -x "$binpath" || fail

check_pr_status() {
  cluster_state=$($binpath read -format=json sys/replication/performance/status | jq -r '.data.state')

  if [[ "${REPLICATION_MODE}" == "primary" ]]; then
    connection_status=$($binpath read -format=json sys/replication/performance/status | jq -r '.data.secondaries[0].connection_status')
  else
    connection_status=$($binpath read -format=json sys/replication/performance/status | jq -r '.data.primaries[0].connection_status')
  fi

  if [[ "$connection_status" == 'disconnected' ]]; then
    fail
  fi

  if [[ "$cluster_state" == 'idle' ]]; then
    fail
  fi
}

# Retry a few times because it can take some time for replication to sync
retry 5 check_pr_status
echo $($binpath read -format=json sys/replication/performance/status)
