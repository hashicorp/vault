# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  aws_mount                  = "aws"     # aws engine
  aws_role                   = "test-role"
  aws_region            = var.aws_test_region
  aws_access_key_id     = var.aws_test_access_key_id
  aws_access_secret_key = var.aws_test_access_secret_key

  // Output
  aws_output = {
    mount                  = local.aws_mount
    role                   = local.aws_role
    region                 = local.aws_region
    test_access_key_id     = local.aws_access_key_id
    test_access_secret_key = local.aws_access_secret_key
  }
}

output "aws_engine" {
  value = local.aws_output
}

# Enable aws secrets engine
resource "enos_remote_exec" "secrets_enable_aws_secret" {
  environment = {
    ENGINE            = local.aws_mount
    MOUNT             = local.aws_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/secrets-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Enable kv secrets engine
resource "enos_remote_exec" "aws_generate_creds" {
  depends_on = [enos_remote_exec.secrets_enable_aws_secret]
  for_each   = var.hosts
  environment = {
    AWS_REGION            = var.aws_test_region
    AWS_ACCESS_KEY_ID     = var.aws_test_access_key_id
    AWS_SECRET_ACCESS_KEY = var.aws_test_access_secret_key
    AWS_ROLE              = local.aws_role
    MOUNT                 = local.aws_mount
    VAULT_ADDR            = var.vault_addr
    VAULT_TOKEN           = var.vault_root_token
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/aws-generate-roles.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
