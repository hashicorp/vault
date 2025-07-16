# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "create_aws_secrets_engine" {
  type        = bool
  description = <<-EOF
    Whether or not we'll verify the AWS secrets engine. Due to the various security requirements in
    Doormat managed AWS accounts, our implementation of the verification requires us to use a
    an external 'DemoUser' role and associated policy in order to create additional users. This is
    configured in vault_ci and vault_enterprise_ci but does not exist in all AWS accounts. As such,
    it's disabled by default.
    See: https://github.com/hashicorp/honeybee-templates/blob/main/templates/iam_policy/DemoUser.yaml
  EOF
  default     = false
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

output "state" {
  value = {
    auth     = local.auth_output
    identity = local.identity_output
    kv       = local.kv_output
    pki      = local.pki_output
    aws      = local.aws_state
  }
}
