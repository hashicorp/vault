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

resource "enos_remote_exec" "configure_pr_secondary" {
  content = templatefile("${path.module}/templates/configure-vault-pr-secondary.sh", {
    vault_cluster_addr      = "${var.secondary_leader_private_ip}:${var.vault_cluster_addr_port}"
    vault_install_dir       = var.vault_install_dir
    vault_local_binary_path = "${var.vault_install_dir}/vault"
    vault_token             = var.vault_root_token
    wrapping_token          = var.wrapping_token
  })

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
    }
  }
}
