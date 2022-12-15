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
  description = "Vault cluster follower Public IP addresses"
}

variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "vault_unseal_keys" {}

locals {
  followers      = toset([for idx in range(var.vault_instance_count - 1) : tostring(idx)])
  vault_bin_path = "${var.vault_install_dir}/vault"
}

resource "enos_remote_exec" "wait_till_sealed" {
  for_each = {
    for idx, follower in local.followers : idx => follower
  }
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    vault_install_dir = var.vault_install_dir
  }

  scripts = ["${path.module}/scripts/wait-till-sealed.sh"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}

resource "enos_remote_exec" "restart_followers" {
  depends_on = [enos_remote_exec.wait_till_sealed]
  for_each = {
    for idx, follower in local.followers : idx => follower
  }

  inline = ["sudo systemctl restart vault"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}

resource "enos_remote_exec" "unseal_followers" {
  depends_on = [enos_remote_exec.restart_followers]
  for_each = {
    for idx, follower in local.followers : idx => follower
    if var.vault_seal_type == "shamir"
  }

  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    vault_install_dir = var.vault_install_dir
    unseal_keys       = join(",", var.vault_unseal_keys)
  }

  scripts = ["${path.module}/scripts/unseal-node.sh"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}

resource "enos_remote_exec" "unseal_followers_again" {
  depends_on = [enos_remote_exec.unseal_followers]
  for_each = {
    for idx, follower in local.followers : idx => follower
    if var.vault_seal_type == "shamir"
  }

  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    vault_install_dir = var.vault_install_dir
    unseal_keys       = join(",", var.vault_unseal_keys)
  }

  scripts = ["${path.module}/scripts/unseal-node.sh"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}
