#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


function fail() {
	echo "$1" 1>&2
	exit 1
}

count=0
retries=20
while :; do
    leader_address=$(curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/v1/sys/leader" | jq '.leader_address' | sed 's/\"//g')
    state=$(curl --header "X-Vault-Token: $VAULT_TOKEN" "$leader_address/v1/sys/metrics"  | jq -r '.Gauges[] | select(.Name == "vault.core.replication.write_undo_logs")')
    target_undo_logs_status="$(jq -r '.Value' <<< "$state")"

    if [ "$target_undo_logs_status" == "1" ]; then
        exit 0
    fi

    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
        echo "$state"
        sleep "$wait"
    else
        fail "Undo_logs did not get into the correct status"
    fi
done
