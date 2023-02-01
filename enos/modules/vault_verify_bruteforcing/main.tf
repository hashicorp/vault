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

variable "node_public_ips" {
  type        = list(string)
  description = "Vault cluster node Public IP address"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  followers      = toset([for idx in range(var.vault_instance_count - 1) : tostring(idx)])
  vault_bin_path = "${var.vault_install_dir}/vault"
}

resource "enos_remote_exec" "verify-bruteforcing" {
  for_each = {
    for idx, follower in local.followers : idx => follower
  }

  environment = {
    VAULT_TOKEN       = var.vault_root_token
    VAULT_ADDR        = "http://localhost:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/verify-bruteforcing.sh")]

  transport = {
    ssh = {
      host = element(var.node_public_ips, each.key)
    }
  }
}