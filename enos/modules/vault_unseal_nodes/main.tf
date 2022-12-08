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

variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "vault_unseal_keys" {}

# variable "vault_unseal_keys" {
#   type        = list(string)
#   description = "Vault cluster unseal keys"
#   default     = null
# }

locals {
  followers      = toset([for idx in range(var.vault_instance_count - 1) : tostring(idx)])
  vault_bin_path = "${var.vault_install_dir}/vault"
  # unseal_keys    = flatten(var.vault_unseal_keys)
}

# # wait for 120s before unsealing the nodes
# resource "time_sleep" "wait_120_seconds" {
#   create_duration = "120s"
# }

resource "enos_vault_unseal" "node" {
  # depends_on = [time_sleep.wait_120_seconds]
  for_each = {
    for idx, follower in local.followers : idx => follower
    if var.vault_seal_type == "shamir"
  }

  bin_path    = local.vault_bin_path
  vault_addr  = "http://localhost:8200"
  seal_type   = var.vault_seal_type
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}
