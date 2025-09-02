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
  description = "The vault cluster instances that were created"
}

variable "ip_version" {
  type        = number
  description = "The IP version to use for the Vault TCP listeners"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_cluster_addr_port" {
  description = "The Raft cluster address port"
  type        = string
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  cluster_addrs = {
    4 : { for k, v in var.hosts : k => "${v.private_ip}:${var.vault_cluster_addr_port}" },
    6 : { for k, v in var.hosts : k => "[${v.ipv6}]:${var.vault_cluster_addr_port}" },
  }
}

resource "enos_remote_exec" "verify_raft_auto_join_voter" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR              = var.vault_addr
    VAULT_CLUSTER_ADDR      = local.cluster_addrs[var.ip_version][each.key]
    VAULT_INSTALL_DIR       = var.vault_install_dir
    VAULT_LOCAL_BINARY_PATH = "${var.vault_install_dir}/vault"
    VAULT_TOKEN             = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/verify-raft-auto-join-voter.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
