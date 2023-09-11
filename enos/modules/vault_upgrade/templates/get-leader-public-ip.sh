#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

binpath=${vault_install_dir}/vault
export VAULT_ADDR="http://localhost:8200"

instances='${vault_instances}'

# Find the leader
leader_address=$($binpath status -format json | jq '.leader_address | rtrimstr(":8200") | ltrimstr("http://")')

# Get the public ip address of the leader
leader_public=$(jq ".[] | select(.private_ip==$leader_address) | .public_ip" <<< "$instances")
echo "$leader_public" | sed 's/\"//g'
