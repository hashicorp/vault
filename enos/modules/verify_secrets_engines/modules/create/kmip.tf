# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

module "create_kmip_secret_engine" {
  depends_on = [
    enos_remote_exec.policy_write_kv_writer,
  ]
  count  = var.kmip_enabled ? 1 : 0
  source = "./kmip"

  integration_host_state = var.integration_host_state
  ip_version             = var.ip_version
  leader_host            = var.leader_host
  ports                  = var.ports
  vault_addr             = var.vault_addr
  vault_edition          = var.vault_edition
  vault_root_token       = var.vault_root_token
  vault_install_dir      = var.vault_install_dir
}

locals {
  kmip_output = var.kmip_enabled ? module.create_kmip_secret_engine[0].kmip : null
}

output "kmip" {
  value = local.kmip_output
}
