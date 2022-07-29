#!/bin/bash

set -eux

binpath=${vault_install_dir}/vault
export VAULT_ADDR="http://localhost:8200"

if "$binpath" status | grep "HA Mode" | grep -ev ${upgrade_target}; then
  exit 0
fi

sudo systemctl restart vault

until "$binpath" status; do
  sleep 1s
done
