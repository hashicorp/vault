# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

module "create_aws_secrets_engine" {
  count  = var.create_aws_secrets_engine ? 1 : 0
  source = "./aws"

  hosts             = var.hosts
  leader_host       = var.leader_host
  vault_addr        = var.vault_addr
  vault_root_token  = var.vault_root_token
  vault_install_dir = var.vault_install_dir
}

locals {
  aws_state = var.create_aws_secrets_engine ? module.create_aws_secrets_engine[0].state : null
}

output "aws" {
  value = local.aws_state
}
