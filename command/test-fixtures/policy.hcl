# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

path "secret/foo" {
  policy = "write"
}

path "secret/bar/*" {
  capabilities = ["create", "read", "update"]
}
