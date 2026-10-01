# Copyright IBM Corp. 2026
# SPDX-License-Identifier: BUSL-1.1

# Auditor-style persona: can browse list and detail views but cannot create, update, or delete.

path "sys/mounts" {
  capabilities = ["read", "list"]
}

path "sys/mounts/*" {
  capabilities = ["read", "list"]
}

path "sys/internal/ui/mounts/*" {
  capabilities = ["read"]
}

# Secrets engines are only listed when the token can reach them.
path "readonly-kv/*" {
  capabilities = ["read", "list"]
}
