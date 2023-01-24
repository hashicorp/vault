#!/bin/bash

undo_logs_status="${VAULT_UNDO_LOGS_STATUS}"

function fail() {
	echo "$1" 1>&2
	exit 1
}

count=0
retries=7
while :; do
    state=$(curl --header "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/v1/sys/metrics"  | jq -r '.Gauges[] | select(.Name == "vault.core.replication.write_undo_logs")')
    target_undo_logs_status="$(jq -r '.Value' <<< "$state")"

    if [ "$undo_logs_status" = "$target_undo_logs_status" ]; then
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
