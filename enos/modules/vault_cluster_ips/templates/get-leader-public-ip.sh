#!/bin/bash

set -e

binpath=${vault_install_dir}/vault
export VAULT_ADDR="http://localhost:8200"

instances='${vault_instances}'

# Find the leader
leader_address=$($binpath status -format json | jq '.leader_address | rtrimstr(":8200") | ltrimstr("http://")')
# leader_address=$("$binpath" operator raft list-peers |grep leader |awk '{print $2}' |awk -F":" '{print $1}')

# Get the public ip address of the leader
leader_public_ip=$(jq ".[] | select(.private_ip==$leader_address) | .public_ip" <<< "$instances")
echo "$leader_public_ip" | sed 's/\"//g'
