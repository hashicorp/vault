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

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

locals {
  vault_bin_path = "${var.vault_install_dir}/vault"
}

resource "enos_remote_exec" "verify_kv_on_node" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/verify-data.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
