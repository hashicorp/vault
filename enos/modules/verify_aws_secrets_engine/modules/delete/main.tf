# Copyright IBM Corp. 2016, 2026
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

# Disable AWS secrets engine
resource "enos_remote_exec" "disable_aws_secrets" {
  environment = {
    MOUNT             = var.create_state.mount
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  inline = ["${var.vault_install_dir}/vault secrets disable ${var.create_state.mount}"]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
