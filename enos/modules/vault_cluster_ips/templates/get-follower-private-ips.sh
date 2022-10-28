#!/bin/bash

set -e

binpath=${vault_install_dir}/vault
export VAULT_ADDR="http://localhost:8200"

instances='${vault_instances}'

# Find the leader
# leader_address=$("$binpath" operator raft list-peers |grep leader |awk '{print $2}' |awk -F":" '{print $1}')
leader_address=$($binpath status -format json | jq '.leader_address | rtrimstr(":8200") | ltrimstr("http://")')

# Get the private ip addresses of the followers
follower_private_ips=$(jq ".[] | select(.private_ip!=$leader_address) | .private_ip" <<< "$instances")

echo "$follower_private_ips" | sed 's/\"//g' | tr '\n' ' '
