# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

module "create_ldap_secret_engine" {
  depends_on = [
    enos_remote_exec.policy_write_kv_writer,
  ]
  count  = var.ldap_enabled ? 1 : 0
  source = "./ldap"

  integration_host_state = var.integration_host_state
  ip_version             = var.ip_version
  leader_host            = var.leader_host
  ports                  = var.ports
  vault_addr             = var.vault_addr
  vault_root_token       = var.vault_root_token
  vault_install_dir      = var.vault_install_dir
}

locals {
  ldap_output = var.ldap_enabled ? module.create_ldap_secret_engine[0].ldap : null
}

output "ldap" {
  value = local.ldap_output
}
