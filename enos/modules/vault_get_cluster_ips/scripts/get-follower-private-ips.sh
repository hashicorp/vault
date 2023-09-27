#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

function fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_INSTANCE_PRIVATE_IPS" ]] && fail "VAULT_INSTANCE_PRIVATE_IPS env variable has not been set"
[[ -z "$VAULT_LEADER_PRIVATE_IP" ]] && fail "VAULT_LEADER_PRIVATE_IP env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "Unable to locate vault binary at $binpath"

count=0
retries=5
while :; do
  # Vault >= 1.10.x has the operator members. If we have that then we'll use it.
  if $binpath operator -h 2>&1 | grep members &> /dev/null; then
    # Get the folllowers that are part of our private ips.
    if followers=$($binpath operator members -format json | jq --argjson expected "$VAULT_INSTANCE_PRIVATE_IPS" -c '.Nodes | map(select(any(.; .active_node==false)) | .api_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")) as $followers | $expected - ($expected - $followers)'); then
      # Make sure that we got all the followers
      if jq --argjson expected "$VAULT_INSTANCE_PRIVATE_IPS" --argjson followers "$followers" -ne '$expected | length as $el | $followers | length as $fl | $fl == $el-1' > /dev/null; then
        echo "$followers"
        exit 0
      fi
    fi
  else
    # We're using an old version of vault so we'll just return ips that don't match the leader.
    # Get the public ip addresses of the followers
    if followers=$(jq --arg ip "$VAULT_LEADER_PRIVATE_IP" -c '. | map(select(.!=$ip))' <<< "$VAULT_INSTANCE_PRIVATE_IPS"); then
      if [[ -n "$followers" ]]; then
        echo "$followers"
        exit 0
      fi
    fi
  fi

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else
    fail "Timed out trying to obtain the cluster followers"
  fi
done
