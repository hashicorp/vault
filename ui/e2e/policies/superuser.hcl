# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

path "*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "sys/leases/lookup" {
    capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
