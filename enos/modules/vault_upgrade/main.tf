# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
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

module "get_ip_addresses" {
  source = "../vault_get_cluster_ips"

  depends_on = [enos_bundle_install.upgrade_vault_binary]

  hosts             = var.hosts
  ip_version        = var.ip_version
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
  vault_root_token  = var.vault_root_token
}

resource "enos_remote_exec" "restart_followers" {
  for_each = module.get_ip_addresses.follower_hosts

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/restart-vault.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_vault_unseal" "followers" {
  for_each = {
    for idx, host in module.get_ip_addresses.follower_hosts : idx => host
    if var.vault_seal_type == "shamir"
  }
  depends_on = [enos_remote_exec.restart_followers]

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
  source = "../vault_verify_unsealed"
  depends_on = [
    enos_remote_exec.restart_followers,
    enos_vault_unseal.followers,
  ]

  hosts             = module.get_ip_addresses.follower_hosts
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
}

resource "enos_remote_exec" "restart_leader" {
  depends_on = [module.wait_for_followers_unsealed]

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/restart-vault.sh")]

  transport = {
    ssh = {
      host = module.get_ip_addresses.leader_public_ip
    }
  }
}

resource "enos_vault_unseal" "leader" {
  count      = var.vault_seal_type == "shamir" ? 1 : 0
  depends_on = [enos_remote_exec.restart_leader]

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
