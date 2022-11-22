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

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "follower_public_ips" {
  type        = list(string)
  description = "Vault primary cluster follower Public IP addresses"
}

variable "vault_unseal_keys" {
  type        = list(string)
  description = "Vault cluster unseal keys"
}

variable "unseal_method" {
  type        = string
  description = "The vault cluster unseal method"
}

locals {
  followers      = toset([for idx in range(var.vault_instance_count - 1) : tostring(idx)])
  vault_bin_path = "${var.vault_install_dir}/vault"
}

# wait for 60s before unsealing the nodes
# resource "time_sleep" "wait_60_seconds" {
#   create_duration = "60s"
# }

resource "enos_vault_unseal" "node" {
  # depends_on = [time_sleep.wait_60_seconds]

  for_each = local.followers

  bin_path    = local.vault_bin_path
  vault_addr  = "http://localhost:8200"
  seal_type   = var.unseal_method
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}
