#!/usr/bin/env bash

set -e

binpath=${vault_install_dir}/vault

instance_ips=${vault_instance_private_ips}

function fail() {
	echo "$1" 1>&2
	exit 1
}

count=0
retries=7
while :; do
    # Find the leader private IP address
    leader_private_ip=$($binpath status -format json | jq '.leader_address | rtrimstr(":8200") | ltrimstr("http://")')
    ip_index=$(echo $instance_ips |jq -r --argjson ip $leader_private_ip 'map(select(. == $ip))')

    if [[ "$leader_private_ip" != 'null' ]] && [[ "$ip_index" != '[]' ]]; then
        echo "$leader_private_ip" | sed 's/\"//g'
        exit 0
    fi

    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
        echo "count is $count for ip $leader_private_ip"
        sleep "$wait"
    else
        fail "leader IP address was not found"
    fi
done
