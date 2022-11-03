#!/bin/bash

set -e

binpath=${vault_install_dir}/vault
export VAULT_ADDR="http://localhost:8200"

instances='${vault_instances}'

# Find the leader private IP address
leader_private_ip=$($binpath status -format json | jq '.leader_address | rtrimstr(":8200") | ltrimstr("http://")')

# Get the private IP address of the leader
echo "$leader_private_ip" | sed 's/\"//g'
