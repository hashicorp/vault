# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "secondary_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The secondary cluster leader host"
}

variable "replication_type" {
  type        = string
  description = "The type of replication to perform"

  validation {
    condition     = contains(["dr", "performance"], var.replication_type)
    error_message = "The replication_type must be either dr or performance"
  }
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

resource "enos_remote_exec" "enable_replication" {
  environment = {
    VAULT_ADDR  = var.vault_addr
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault write sys/replication/${var.replication_type}/secondary/enable token=${var.wrapping_token}"]

  transport = {
    ssh = {
      host = var.secondary_leader_host.public_ip
    }
  }
}

// Wait for our primary host to be the "leader", which means it's running and all "setup" tasks
// have been completed. We'll have to unseal our follower nodes after this has occurred.
module "wait_for_leader" {
  source = "../vault_wait_for_leader"

  depends_on = [
    enos_remote_exec.enable_replication
  ]

  hosts             = { "0" : var.secondary_leader_host }
  ip_version        = var.ip_version
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
  vault_root_token  = var.vault_root_token
}

// Ensure that our leader is ready to for us to unseal follower nodes.
resource "enos_remote_exec" "wait_for_leader_ready" {
  depends_on = [
    module.wait_for_leader,
  ]

  environment = {
    REPLICATION_TYPE  = var.replication_type
    RETRY_INTERVAL    = 3  // seconds
    TIMEOUT_SECONDS   = 60 // seconds
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/wait-for-leader-ready.sh")]

  transport = {
    ssh = {
      host = var.secondary_leader_host.public_ip
    }
  }
}
