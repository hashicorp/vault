# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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
  description = "The Vault cluster instances that were created"
}

variable "leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })

  description = "Vault cluster leader host"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
  default     = null
}

# We use this module to verify write data in all Enos scenarios.  Since we cannot use
# Vault token to authenticate to secondary clusters in replication scenario we add a regular user
# here to keep the authentication method and module verification consistent between all scenarios
resource "enos_remote_exec" "smoke-enable-secrets-kv" {
  # Only enable the secrets engine on the leader node
  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/smoke-enable-secrets-kv.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Verify that we can enable the k/v secrets engine and write data to it.
resource "enos_remote_exec" "smoke-write-test-data" {
  depends_on = [enos_remote_exec.smoke-enable-secrets-kv]
  for_each   = var.hosts

  environment = {
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
    TEST_KEY          = "smoke${each.key}"
    TEST_VALUE        = "fire"
  }

  scripts = [abspath("${path.module}/scripts/smoke-write-test-data.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
