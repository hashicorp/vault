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

variable "primary_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The primary cluster leader host"
}

variable "secondary_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The secondary cluster leader host"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "wrapping_token" {
  type        = string
  description = "The wrapping token created on primary cluster"
  default     = null
}

locals {
  primary_leader_addr          = var.ip_version == 6 ? var.primary_leader_host.ipv6 : var.primary_leader_host.private_ip
  secondary_leader_addr        = var.ip_version == 6 ? var.secondary_leader_host.ipv6 : var.secondary_leader_host.private_ip
  primary_replication_status   = jsondecode(enos_remote_exec.verify_replication_status_on_primary.stdout)
  secondary_replication_status = jsondecode(enos_remote_exec.verify_replication_status_on_secondary.stdout)
}

resource "enos_remote_exec" "verify_replication_status_on_primary" {
  environment = {
    IP_VERSION            = var.ip_version
    PRIMARY_LEADER_ADDR   = local.primary_leader_addr
    SECONDARY_LEADER_ADDR = local.secondary_leader_addr
    VAULT_ADDR            = var.vault_addr
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/verify-replication-status.sh")]

  transport = {
    ssh = {
      host = var.primary_leader_host.public_ip
    }
  }
}

resource "enos_remote_exec" "verify_replication_status_on_secondary" {
  environment = {
    IP_VERSION            = var.ip_version
    PRIMARY_LEADER_ADDR   = local.primary_leader_addr
    SECONDARY_LEADER_ADDR = local.secondary_leader_addr
    VAULT_ADDR            = var.vault_addr
    VAULT_INSTALL_DIR     = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/verify-replication-status.sh")]

  transport = {
    ssh = {
      host = var.secondary_leader_host.public_ip
    }
  }
}

output "primary_replication_status" {
  value = local.primary_replication_status
}

output "known_primary_cluster_addrs" {
  value = local.secondary_replication_status.data.known_primary_cluster_addrs
}

output "secondary_replication_status" {
  value = local.secondary_replication_status
}

output "primary_replication_data_secondaries" {
  value = local.primary_replication_status.data.secondaries
}

output "secondary_replication_data_primaries" {
  value = local.secondary_replication_status.data.primaries
}
