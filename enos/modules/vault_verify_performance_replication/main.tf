# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_cluster_addr_port" {
  description = "The Raft cluster address port"
  type        = string
  default     = "8201"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "primary_leader_public_ip" {
  type        = string
  description = "Vault primary cluster leader Public IP address"
}

variable "primary_leader_private_ip" {
  type        = string
  description = "Vault primary cluster leader Private IP address"
}

variable "secondary_leader_public_ip" {
  type        = string
  description = "Vault secondary cluster leader Public IP address"
}

variable "secondary_leader_private_ip" {
  type        = string
  description = "Vault secondary cluster leader Private IP address"
}

variable "wrapping_token" {
  type        = string
  description = "The wrapping token created on primary cluster"
  default     = null
}

locals {
  primary_replication_status   = jsondecode(enos_remote_exec.verify_replication_status_on_primary.stdout)
  secondary_replication_status = jsondecode(enos_remote_exec.verify_replication_status_on_secondary.stdout)
}

resource "enos_remote_exec" "verify_replication_status_on_primary" {
  environment = {
    VAULT_ADDR               = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR        = var.vault_install_dir
    PRIMARY_LEADER_PRIV_IP   = var.primary_leader_private_ip
    SECONDARY_LEADER_PRIV_IP = var.secondary_leader_private_ip
  }

  scripts = [abspath("${path.module}/scripts/verify-replication-status.sh")]

  transport = {
    ssh = {
      host = var.primary_leader_public_ip
    }
  }
}

resource "enos_remote_exec" "verify_replication_status_on_secondary" {
  environment = {
    VAULT_ADDR               = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR        = var.vault_install_dir
    PRIMARY_LEADER_PRIV_IP   = var.primary_leader_private_ip
    SECONDARY_LEADER_PRIV_IP = var.secondary_leader_private_ip
  }

  scripts = [abspath("${path.module}/scripts/verify-replication-status.sh")]

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
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
