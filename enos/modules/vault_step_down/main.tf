# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "leader_host" {
  type = object({
    private_ip = string
    public_ip  = string
  })

  description = "The vault cluster host that can be expected as a leader"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

resource "enos_remote_exec" "vault_operator_step_down" {
  environment = {
    VAULT_TOKEN       = var.vault_root_token
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/operator-step-down.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
