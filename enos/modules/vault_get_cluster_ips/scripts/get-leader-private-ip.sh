#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

binpath=${VAULT_INSTALL_DIR}/vault
instance_ips=${VAULT_INSTANCE_PRIVATE_IPS}

function fail() {
	echo "$1" 1>&2
	exit 1
}

count=0
retries=5
while :; do
    # Find the leader private IP address
    leader_private_ip=$($binpath status -format json | jq '.leader_address | rtrimstr(":8200") | ltrimstr("http://")')
    match_ip=$(echo "$instance_ips" |jq -r --argjson ip "$leader_private_ip" 'map(select(. == $ip))')

    if [[ "$leader_private_ip" != 'null' ]] && [[ "$match_ip" != '[]' ]]; then
        echo "$leader_private_ip" | sed 's/\"//g'
        exit 0
    fi

    wait=$((5 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
        sleep "$wait"
    else
        fail "leader IP address $leader_private_ip was not found in $instance_ips"
    fi
done
