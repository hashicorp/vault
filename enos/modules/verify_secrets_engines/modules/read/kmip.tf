# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Verifying Vault LDAP Configurations
module "verify_kmip_secret_engine" {
  count  = var.kmip_enabled ? 1 : 0
  source = "./kmip"

  create_state  = var.create_state
  hosts         = var.hosts
  ip_version    = var.ip_version
  vault_addr    = var.vault_addr
  vault_edition = var.vault_edition
}

