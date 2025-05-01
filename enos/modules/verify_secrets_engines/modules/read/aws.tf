# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

module "verify_aws_secrets_engine" {
  count  = var.verify_aws_secrets_engine ? 1 : 0
  source = "./aws"

  create_state            = var.create_state
  vault_addr              = var.vault_addr
  vault_root_token        = var.vault_root_token
  vault_install_dir       = var.vault_install_dir
  verify_aws_engine_creds = var.verify_aws_engine_creds

  hosts = var.hosts
}
