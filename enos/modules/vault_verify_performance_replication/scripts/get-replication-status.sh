#!/usr/bin/env bash

set -e

binpath=${VAULT_INSTALL_DIR}/vault

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
      echo $pr_status
      return "$exit"
    fi
  done

  echo $pr_status
  return 0
}

test -x "$binpath" || exit 1

check_pr_status() {
  pr_status=$($binpath read -format=json sys/replication/performance/status)
  cluster_state=$(echo $pr_status | jq -r '.data.state')

  if [[ "${REPLICATION_MODE}" == "primary" ]]; then
    connection_status=$(echo $pr_status | jq -r '.data.secondaries[0].connection_status')
  else
    connection_status=$(echo $pr_status | jq -r '.data.primaries[0].connection_status')
  fi

  if [[ "$connection_status" == 'disconnected' ]] || [[ "$cluster_state" == 'idle' ]]; then
    return 1
  fi
}

# Retry a few times because it can take some time for replication to sync
retry 5 check_pr_status
