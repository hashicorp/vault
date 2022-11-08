terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "follower_public_ip" {
  type        = string
  description = "Vault cluster follower Public IP addresses"
}

variable "vault_unseal_keys" {
  type        = list(string)
  description = "Vault cluster unseal keys"
}

variable "unseal_method" {
  type        = string
  description = "The vault cluster unseal method"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  vault_bin_path = "${var.vault_install_dir}/vault"
}
resource "enos_vault_unseal" "node" {
  bin_path    = local.vault_bin_path
  vault_addr  = "http://localhost:8200"
  seal_type   = var.unseal_method
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = var.follower_public_ip
    }
  }
}
