# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Verify PKI Certificate
resource "enos_remote_exec" "aws_verify_new_creds" {
  for_each = var.hosts

  environment = {
    AWS_REGION            = "${var.create_state.aws.region}"
    MOUNT                 = "${var.create_state.aws.mount}"
    AWS_USER_NAME         = "${var.create_state.aws.aws_user_name}"
    AWS_ACCESS_KEY_ID     = "${var.create_state.aws.aws_access_key}"
    AWS_SECRET_ACCESS_KEY = "${var.create_state.aws.aws_secret_key}"
    VAULT_AWS_ROLE        = "${var.create_state.aws.vault_aws_role}"
    VAULT_ADDR            = var.vault_addr
    VAULT_TOKEN           = var.vault_root_token
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/aws-verify-new-creds.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

