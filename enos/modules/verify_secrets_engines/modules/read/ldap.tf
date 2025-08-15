# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Verifying Vault LDAP Configurations
module "verify_ldap_secret_engine" {
  count  = var.ldap_enabled ? 1 : 0
  source = "./ldap"

  create_state      = var.create_state
  vault_addr        = var.vault_addr
  vault_root_token  = var.vault_root_token
  vault_install_dir = var.vault_install_dir

  hosts = var.hosts
}

