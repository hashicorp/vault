# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })

  description = "Vault cluster leader host"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
  default     = null
}

locals {
  // Variables
  aws_mount      = "aws"
  vault_aws_role = "enos_test_role"
  my_email       = split("/", data.aws_caller_identity.current.arn)[2]

  // State output
  state = {
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

output "state" {
  value = local.state
}

resource "random_id" "unique_suffix" {
  byte_length = 4
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

# The "DemoUser" policy is a predefined policy created by the security team.
# This policy grants the necessary AWS permissions required for role generation via Vault.
# Reference: https://github.com/hashicorp/honeybee-templates/blob/main/templates/iam_policy/DemoUser.yaml
data "aws_iam_policy" "premade_demo_user_policy" {
  name = "DemoUser"
}

# This role was provisioned by the security team using the repository referenced below.
# This role includes the necessary policies to enable AWS credential generation and rotation via Vault.
# Reference: https://github.com/hashicorp/honeybee-templates/blob/main/templates/iam_role/vault-assumed-role-credentials-demo.yaml
data "aws_iam_role" "premade_demo_assumed_role" {
  name = "vault-assumed-role-credentials-demo"
}

# Creating new test user
resource "aws_iam_user" "aws_enos_test_user" {
  name                 = "demo-${local.my_email}-${random_id.unique_suffix.hex}"
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

  scripts = [abspath("${path.module}/../../../scripts/secrets-enable.sh")]

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
    AWS_REGION            = local.state.region
    ENGINE                = local.aws_mount
    MOUNT                 = local.aws_mount
    AWS_USER_NAME         = local.state.aws_user_name
    AWS_POLICY_ARN        = local.state.aws_policy_arn
    AWS_ROLE_ARN          = local.state.aws_role_arn
    AWS_ACCESS_KEY_ID     = local.state.aws_access_key
    AWS_SECRET_ACCESS_KEY = local.state.aws_secret_key
    VAULT_AWS_ROLE        = local.vault_aws_role
    VAULT_ADDR            = var.vault_addr
    VAULT_TOKEN           = var.vault_root_token
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../../scripts/aws-generate-roles.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
