# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.5.4"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}


variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_artifactory_release" {
  type = object({
    username = string
    token    = string
    url      = string
    sha256   = string
  })
  description = "Vault release version and edition to install from artifactory.hashicorp.engineering"
  default     = null
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_local_artifact_path" {
  type        = string
  description = "The path to a locally built vault artifact to install"
  default     = null
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "vault_unseal_keys" {
  type        = list(string)
  description = "The keys to use to unseal Vault when not using auto-unseal"
  default     = null
}

locals {
  vault_bin_path = "${var.vault_install_dir}/vault"
}

// Upgrade the Vault artifact in-place. With zip bundles we must use the same path of the original
// installation so that we can re-use the systemd unit that enos_vault_start created at
// /etc/systemd/system/vault.service. The path does not matter for package types as the systemd
// unit for the bianry is included and will be installed.
resource "enos_bundle_install" "upgrade_vault_binary" {
  for_each = var.hosts

  destination = var.vault_install_dir
  artifactory = var.vault_artifactory_release
  path        = var.vault_local_artifact_path

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

// We assume that our original Vault cluster used a zip bundle from releases.hashicorp.com and as
// such enos_vault_start will have created a systemd unit for it at /etc/systemd/systemd/vault.service.
// If we're upgrading to a package that contains its own systemd unit we'll need to remove the
// old unit file so that when we restart vault we pick up the new unit that points to the updated
// binary.
resource "enos_remote_exec" "maybe_remove_old_unit_file" {
  for_each   = var.hosts
  depends_on = [enos_bundle_install.upgrade_vault_binary]

  environment = {
    ARTIFACT_NAME = enos_bundle_install.upgrade_vault_binary[each.key].name
  }

  scripts = [abspath("${path.module}/scripts/maybe-remove-old-unit-file.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

module "get_ip_addresses" {
  source = "../vault_get_cluster_ips"

  depends_on = [enos_remote_exec.maybe_remove_old_unit_file]

  hosts             = var.hosts
  ip_version        = var.ip_version
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
  vault_root_token  = var.vault_root_token
}

module "restart_followers" {
  source            = "../restart_vault"
  hosts             = module.get_ip_addresses.follower_hosts
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
}

resource "enos_vault_unseal" "followers" {
  for_each = {
    for idx, host in module.get_ip_addresses.follower_hosts : idx => host
    if var.vault_seal_type == "shamir"
  }
  depends_on = [module.restart_followers]

  bin_path    = local.vault_bin_path
  vault_addr  = var.vault_addr
  seal_type   = var.vault_seal_type
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

module "wait_for_followers_unsealed" {
  source = "../vault_wait_for_cluster_unsealed"
  depends_on = [
    module.restart_followers,
    enos_vault_unseal.followers,
  ]

  hosts             = module.get_ip_addresses.follower_hosts
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
}

module "restart_leader" {
  depends_on        = [module.wait_for_followers_unsealed]
  source            = "../restart_vault"
  hosts             = module.get_ip_addresses.leader_hosts
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
}

resource "enos_vault_unseal" "leader" {
  count      = var.vault_seal_type == "shamir" ? 1 : 0
  depends_on = [module.restart_leader]

  bin_path    = local.vault_bin_path
  vault_addr  = var.vault_addr
  seal_type   = var.vault_seal_type
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = module.get_ip_addresses.leader_public_ip
    }
  }
}
