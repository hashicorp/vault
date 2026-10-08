# Copyright IBM Corp. 2026
# SPDX-License-Identifier: BUSL-1.1

# Auditor-style persona: can browse list and detail views but cannot create, update, or delete.

path "sys/namespaces" {
  capabilities = ["read", "list"]
}

path "sys/namespaces/*" {
  capabilities = ["read", "list"]
}

path "sys/mounts" {
  capabilities = ["read", "list"]
}

path "sys/mounts/*" {
  capabilities = ["read", "list"]
}

path "sys/internal/ui/mounts/*" {
  capabilities = ["read"]
}

path "sys/policies/+" {
  capabilities = ["list"]
}

# Only policies named readable-* can be read, so list views show policies the token cannot read.
path "sys/policies/+/readable-*" {
  capabilities = ["read"]
}

path "identity/*" {
  capabilities = ["read", "list"]
}

# Secrets engines are only listed when the token can reach them.
path "readonly-kv/*" {
  capabilities = ["read", "list"]
}

# LDAP roles can be viewed but not edited, deleted or rotated.
path "readonly-ldap/*" {
  capabilities = ["read", "list"]
}
