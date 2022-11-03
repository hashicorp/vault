terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
    time = {
      source = "hashicorp/time"
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

variable "primary_vault_root_token" {
  type        = string
  description = "The vault root token of primary cluster"
}

variable "secondary_vault_root_token" {
  type        = string
  description = "The vault root token of secondary cluster"
}

variable "wrapping_token" {
  type        = string
  description = "The wrapping token created on primary cluster"
  default     = null
}

locals {
  primary_replication_status = jsondecode(enos_remote_exec.verify_replication_on_primary.stdout)
  secondary_replication_status = jsondecode(enos_remote_exec.verify_replication_on_secondary.stdout)
}
resource "enos_remote_exec" "verify_replication_on_primary" {
  environment = {
    VAULT_ADDR  = "http://127.0.0.1:8200"
    VAULT_TOKEN = var.primary_vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault read -format=json sys/replication/performance/status"]

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

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
}

resource "enos_remote_exec" "verify_replication_on_secondary" {
  # wait 60s before verifying the status on secondary
  depends_on = [time_sleep.wait_60_seconds]
  environment = {
    VAULT_ADDR  = "http://127.0.0.1:8200"
    VAULT_TOKEN = var.secondary_vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault read -format=json sys/replication/performance/status"]

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
    }
  }
}

output "secondary_replication_status" {
  value = local.secondary_replication_status

  precondition {
    condition     = local.secondary_replication_status.data.mode == "secondary" && local.secondary_replication_status.data.state != "idle"
    error_message = "Vault secondary cluster mode must be \"secondary\" and state must not be \"idle\"."
  }
}
