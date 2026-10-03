#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$IP_VERSION" ]] && fail "IP_VERSION env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "Unable to locate vault binary at $binpath"

getFollowerPrivateIPsFromOperatorMembers() {
  if members=$($binpath operator members -format json); then
    if followers=$(echo "$members" | jq -e --argjson expected "$VAULT_PRIVATE_IPS" -c '.Nodes | map(select(any(.; .active_node==false)) | .api_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")) as $followers | $expected - ($expected - $followers)'); then
      # Make sure that we got all the followers
      if jq -e --argjson expected "$VAULT_PRIVATE_IPS" --argjson followers "$followers" -ne '$expected | length as $el | $followers | length as $fl | $fl == $el-1' > /dev/null; then
        echo "$followers"
        return 0
      fi
    fi
  fi

  return 1
}

removeIP() {
  local needle
  local haystack
  needle=$1
  haystack=$2
  if remain=$(jq -e --arg ip "$needle" -c '. | map(select(.!=$ip))' <<< "$haystack"); then
    if [[ -n "$remain" ]]; then
      echo "$remain"
      return 0
    fi
  fi

  return 1
}

count=0
retries=10
while :; do
  case $IP_VERSION in
    4)
      [[ -z "$VAULT_PRIVATE_IPS" ]] && fail "VAULT_PRIVATE_IPS env variable has not been set"
      [[ -z "$VAULT_LEADER_PRIVATE_IP" ]] && fail "VAULT_LEADER_PRIVATE_IP env variable has not been set"

      # Vault >= 1.10.x has the operator members. If we have that then we'll use it.
      if $binpath operator -h 2>&1 | grep members &> /dev/null; then
        if followers=$(getFollowerPrivateIPsFromOperatorMembers); then
          echo "$followers"
          exit 0
        fi
      else
        removeIP "$VAULT_LEADER_PRIVATE_IP" "$VAULT_PRIVATE_IPS"

        return $?
      fi
      ;;
    6)
      echo '[]'
      exit 0
      ;;
    *)
      fail "unknown IP_VERSION: $IP_VERSION"
      ;;
  esac

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else
    fail "Timed out trying to obtain the cluster followers"
  fi
done
