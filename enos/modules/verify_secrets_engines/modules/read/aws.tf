# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Verify PKI Certificate
resource "enos_remote_exec" "aws_verify_roles" {
  for_each = var.hosts

  environment = {
    AWS_REGION    = var.create_state.aws.region
    AWS_ACCESS_KEY_ID     = var.create_state.aws.test_access_key_id
    AWS_SECRET_ACCESS_KEY = var.create_state.aws.test_access_secret_key
    AWS_ROLE              = var.create_state.aws.role
    MOUNT                 = var.create_state.aws.mount
    VAULT_ADDR            = var.vault_addr
    VAULT_TOKEN           = var.vault_root_token
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/aws-verify-roles.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

