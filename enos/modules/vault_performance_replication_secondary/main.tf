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

variable "secondary_leader_public_ip" {
  type        = string
  description = "Vault secondary cluster leader Public IP address"
}

variable "secondary_leader_private_ip" {
  type        = string
  description = "Vault secondary cluster leader Private IP address"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "wrapping_token" {
  type        = string
  description = "The wrapping token created on primary cluster"
}

locals {
  wrapping_token = var.wrapping_token
}

resource "enos_remote_exec" "configure_pr_secondary" {
  environment = {
    VAULT_ADDR  = "http://127.0.0.1:8200"
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault write sys/replication/performance/secondary/enable token=${local.wrapping_token}"]

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
    }
  }
}
