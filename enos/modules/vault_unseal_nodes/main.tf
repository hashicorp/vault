# This module unseals the replication secondary follower nodes
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

# After replication is enabled the secondary follower nodes are expected to be sealed,
# so we wait for the secondary follower nodes to update the seal status
resource "enos_remote_exec" "wait_until_sealed" {
  for_each = {
    for idx, follower in local.followers : idx => follower
  }
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = ["${path.module}/scripts/wait-until-sealed.sh"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}

# The follower nodes on secondary replication cluster incorrectly report
# unseal progress 2/3 (Issue: https://hashicorp.atlassian.net/browse/VAULT-12309),
# so we restart the followers to clear the status and to autounseal incase of awskms seal type
resource "enos_remote_exec" "restart_followers" {
  depends_on = [enos_remote_exec.wait_until_sealed]
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

# We cannot use the vault_unseal resouce due to the known issue
# (https://hashicorp.atlassian.net/browse/VAULT-12311). We use a custom
# script to allow retry for unsealing the secondary followers
resource "enos_remote_exec" "unseal_followers" {
  depends_on = [enos_remote_exec.restart_followers]
  # The unseal keys are required only for seal_type shamir
  for_each = {
    for idx, follower in local.followers : idx => follower
    if var.vault_seal_type == "shamir"
  }

  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
    UNSEAL_KEYS       = join(",", var.vault_unseal_keys)
  }

  scripts = ["${path.module}/scripts/unseal-node.sh"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}

# This is a second attempt needed to unseal the secondary followers
# using a custom script due to get past the known issue
# (https://hashicorp.atlassian.net/browse/VAULT-12311)
resource "enos_remote_exec" "unseal_followers_again" {
  depends_on = [enos_remote_exec.unseal_followers]
  for_each = {
    for idx, follower in local.followers : idx => follower
    if var.vault_seal_type == "shamir"
  }

  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
    UNSEAL_KEYS       = join(",", var.vault_unseal_keys)
  }

  scripts = ["${path.module}/scripts/unseal-node.sh"]

  transport = {
    ssh = {
      host = element(var.follower_public_ips, each.key)
    }
  }
}
