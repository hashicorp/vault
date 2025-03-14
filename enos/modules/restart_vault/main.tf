# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
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
  description = "The vault hosts"
}

variable "vault_addr" {
  type        = string
  description = "The local vault api address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the vault binary is installed"
}


resource "enos_remote_exec" "restart" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/restart-vault.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

