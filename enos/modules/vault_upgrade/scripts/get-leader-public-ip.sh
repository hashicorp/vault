#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

binpath=${VAULT_INSTALL_DIR}/vault
export VAULT_ADDR="http://localhost:8200"

instances=${VAULT_INSTANCES}

# Find the leader
leader_address=$($binpath status -format json | jq '.leader_address | scan("[0-9]+.[0-9]+.[0-9]+.[0-9]+")')

# Get the public ip address of the leader
leader_public=$(jq ".[] | select(.private_ip==$leader_address) | .public_ip" <<< "$instances")
#shellcheck disable=SC2001
echo "$leader_public" | sed 's/\"//g'
