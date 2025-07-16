# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# This module unseals the replication secondary follower nodes
terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster hosts to unseal"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "vault_unseal_keys" {}

locals {
  vault_bin_path = "${var.vault_install_dir}/vault"
}

# After replication is enabled the secondary follower nodes are expected to be sealed,
# so we wait for the secondary follower nodes to update the seal status
resource "enos_remote_exec" "wait_until_sealed" {
  for_each = var.hosts
  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/wait-until-sealed.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# The follower nodes on secondary replication cluster incorrectly report
# unseal progress 2/3 (Issue: https://hashicorp.atlassian.net/browse/VAULT-12309),
# so we restart the followers to allow them to auto-unseal
resource "enos_remote_exec" "restart_followers" {
  depends_on = [enos_remote_exec.wait_until_sealed]
  for_each = {
    for idx, host in var.hosts : idx => host
    if var.vault_seal_type != "shamir"
  }

  inline = ["sudo systemctl restart vault"]

  transport = {
    ssh = {
      host = each.value.public_ip
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
    for idx, host in var.hosts : idx => host
    if var.vault_seal_type == "shamir"
  }

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    UNSEAL_KEYS       = join(",", var.vault_unseal_keys)
  }

  scripts = [abspath("${path.module}/scripts/unseal-node.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# This is a second attempt needed to unseal the secondary followers
# using a custom script due to get past the known issue
# (https://hashicorp.atlassian.net/browse/VAULT-12311)
resource "enos_remote_exec" "unseal_followers_again" {
  depends_on = [enos_remote_exec.unseal_followers]
  for_each = {
    for idx, host in var.hosts : idx => host
    if var.vault_seal_type == "shamir"
  }

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    UNSEAL_KEYS       = join(",", var.vault_unseal_keys)
  }

  scripts = [abspath("${path.module}/scripts/unseal-node.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
