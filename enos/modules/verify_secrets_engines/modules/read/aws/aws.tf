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

variable "create_state" {
  description = "The state of the secrets engines from the 'create' module"
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

variable "verify_aws_engine_creds" {
  type = bool
}

# Verify AWS Engine
resource "enos_remote_exec" "aws_verify_new_creds" {
  for_each = var.hosts

  environment = {
    AWS_REGION              = "${var.create_state.aws.region}"
    MOUNT                   = "${var.create_state.aws.mount}"
    AWS_USER_NAME           = "${var.create_state.aws.aws_user_name}"
    AWS_ACCESS_KEY_ID       = "${var.create_state.aws.aws_access_key}"
    AWS_SECRET_ACCESS_KEY   = "${var.create_state.aws.aws_secret_key}"
    VAULT_AWS_ROLE          = "${var.create_state.aws.vault_aws_role}"
    VAULT_ADDR              = var.vault_addr
    VAULT_TOKEN             = var.vault_root_token
    VAULT_INSTALL_DIR       = var.vault_install_dir
    VERIFY_AWS_ENGINE_CERTS = var.verify_aws_engine_creds
  }

  scripts = [abspath("${path.module}/../../../scripts/aws-verify-new-creds.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
