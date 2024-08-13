# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "secondary_leader_public_ip" {
  type        = string
  description = "Vault secondary cluster leader Public IP address"
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
  description = "The vault root token"
}

variable "wrapping_token" {
  type        = string
  description = "The wrapping token created on primary cluster"
}

resource "enos_remote_exec" "configure_pr_secondary" {
  environment = {
    VAULT_ADDR  = var.vault_addr
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault write sys/replication/performance/secondary/enable token=${var.wrapping_token}"]

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
    }
  }
}
