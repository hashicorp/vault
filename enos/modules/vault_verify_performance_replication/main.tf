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
  primary_replication_status   = jsondecode(enos_remote_exec.verify_replication_on_primary.stdout)
  secondary_replication_status = jsondecode(enos_remote_exec.verify_replication_on_secondary.stdout)
}

resource "enos_remote_exec" "verify_replication_on_primary" {
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = ["${path.module}/scripts/verify-performance-replication.sh"]

  transport = {
    ssh = {
      host = var.primary_leader_public_ip
    }
  }
}

output "primary_replication_status" {
  value = local.primary_replication_status

  precondition {
    condition     = local.primary_replication_status.data.mode == "primary" && local.primary_replication_status.data.state != "idle"
    error_message = "Vault primary cluster mode must be \"primary\" and state must not be \"idle\"."
  }
}

resource "enos_remote_exec" "verify_replication_on_secondary" {
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = ["${path.module}/scripts/verify-performance-replication.sh"]

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
    }
  }
}

output "known_primary_cluster_addrs" {
  value = local.secondary_replication_status.data.known_primary_cluster_addrs

  precondition {
    condition     = contains(local.secondary_replication_status.data.known_primary_cluster_addrs, "https://${var.primary_leader_private_ip}:8201")
    error_message = "Vault secondary cluster known_primary_cluster_addrs must include ${var.primary_leader_private_ip}."
  }
}

output "secondary_replication_status" {
  value = local.secondary_replication_status

  precondition {
    condition     = local.secondary_replication_status.data.mode == "secondary" && local.secondary_replication_status.data.state != "idle"
    error_message = "Vault secondary cluster mode must be \"secondary\" and state must not be \"idle\"."
  }
}

output "primary_replication_data_secondaries" {
  value = local.primary_replication_status.data.secondaries

  # The secondaries connection_status should be "connected"
  precondition {
    condition     = local.primary_replication_status.data.secondaries[0].connection_status == "connected"
    error_message = "connection status to primaries must be \"connected\"."
  }

  # The secondaries cluster address must have the secondary leader address
  precondition {
    condition     = local.primary_replication_status.data.secondaries[0].cluster_address == "https://${var.secondary_leader_private_ip}:8201"
    error_message = "Vault secondaries cluster_address must be with ${var.secondary_leader_private_ip}."
  }
}

output "secondary_replication_data_primaries" {
  value = local.secondary_replication_status.data.primaries

  # The primaries connection_status should be "connected"
  precondition {
    condition     = local.secondary_replication_status.data.primaries[0].connection_status == "connected"
    error_message = "connection status to primaries must be \"connected\"."
  }

  # The primaries cluster address must have the primary leader address
  precondition {
    condition     = local.secondary_replication_status.data.primaries[0].cluster_address == "https://${var.primary_leader_private_ip}:8201"
    error_message = "Vault primaries cluster_address must be ${var.primary_leader_private_ip}."
  }
}
