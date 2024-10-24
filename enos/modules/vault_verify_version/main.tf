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

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_build_date" {
  type        = string
  description = "The Vault artifact build date"
  default     = null
}

variable "vault_edition" {
  type        = string
  description = "The Vault product edition"
  default     = null
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_product_version" {
  type        = string
  description = "The Vault product version"
  default     = null
}

variable "vault_revision" {
  type        = string
  description = "The Vault product revision"
  default     = null
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
  default     = null
}

resource "enos_remote_exec" "verify_cli_version" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR        = var.vault_addr,
    VAULT_BUILD_DATE  = var.vault_build_date,
    VAULT_EDITION     = var.vault_edition,
    VAULT_INSTALL_DIR = var.vault_install_dir,
    VAULT_REVISION    = var.vault_revision,
    VAULT_TOKEN       = var.vault_root_token,
    VAULT_VERSION     = var.vault_product_version,
  }

  scripts = [abspath("${path.module}/scripts/verify-cli-version.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "verify_cluster_version" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR       = var.vault_addr,
    VAULT_BUILD_DATE = var.vault_build_date,
    VAULT_TOKEN      = var.vault_root_token,
    VAULT_VERSION    = var.vault_product_version,
  }

  scripts = [abspath("${path.module}/scripts/verify-cluster-version.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
