# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  aws_mount      = "aws"
  vault_aws_role = "enos_test_role"
  my_email       = split("/", data.aws_caller_identity.current.arn)[2]

  // Output
  aws_output = {
    aws_role       = data.aws_iam_role.premade_demo_assumed_role.name
    aws_role_arn   = data.aws_iam_role.premade_demo_assumed_role.arn
    aws_policy_arn = data.aws_iam_policy.premade_demo_user_policy.arn
    aws_user_name  = aws_iam_user.aws_enos_test_user.name
    aws_access_key = aws_iam_access_key.aws_enos_test_user.id
    aws_secret_key = aws_iam_access_key.aws_enos_test_user.secret
    mount          = local.aws_mount
    region         = data.aws_region.current.name
    vault_aws_role = local.vault_aws_role
  }
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

# Using Pre-made policy and role
data "aws_iam_policy" "premade_demo_user_policy" {
  name = "DemoUser"
}

data "aws_iam_role" "premade_demo_assumed_role" {
  name = "vault-assumed-role-credentials-demo"
}

# Creating new test user
resource "aws_iam_user" "aws_enos_test_user" {
  name                 = "demo-${local.my_email}-${formatdate("HHmmss", timestamp())}"
  permissions_boundary = data.aws_iam_policy.premade_demo_user_policy.arn
  force_destroy        = true
}

resource "aws_iam_user_policy_attachment" "aws_enos_test_user" {
  user       = aws_iam_user.aws_enos_test_user.name
  policy_arn = data.aws_iam_policy.premade_demo_user_policy.arn
}

resource "aws_iam_access_key" "aws_enos_test_user" {
  user = aws_iam_user.aws_enos_test_user.name
  lifecycle {
    prevent_destroy = false
  }
}

# Enable AWS secrets engine
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

# Generate AWS Credentials
resource "enos_remote_exec" "aws_generate_roles" {
  depends_on = [enos_remote_exec.secrets_enable_aws_secret]
  for_each   = var.hosts

  environment = {
    AWS_REGION            = local.aws_output.region
    ENGINE                = local.aws_mount
    MOUNT                 = local.aws_mount
    AWS_USER_NAME         = local.aws_output.aws_user_name
    AWS_POLICY_ARN        = local.aws_output.aws_policy_arn
    AWS_ROLE_ARN          = local.aws_output.aws_role_arn
    AWS_ACCESS_KEY_ID     = local.aws_output.aws_access_key
    AWS_SECRET_ACCESS_KEY = local.aws_output.aws_secret_key
    VAULT_AWS_ROLE        = local.vault_aws_role
    VAULT_ADDR            = var.vault_addr
    VAULT_TOKEN           = var.vault_root_token
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/aws-generate-roles.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
