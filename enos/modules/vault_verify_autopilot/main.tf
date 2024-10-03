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
  description = "The vault cluster instances that were created"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_autopilot_upgrade_status" {
  type        = string
  description = "The autopilot upgrade expected status"
}

variable "vault_autopilot_upgrade_version" {
  type        = string
  description = "The Vault upgraded version"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

resource "enos_remote_exec" "smoke-verify-autopilot" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR                      = var.vault_addr
    VAULT_INSTALL_DIR               = var.vault_install_dir,
    VAULT_TOKEN                     = var.vault_root_token,
    VAULT_AUTOPILOT_UPGRADE_STATUS  = var.vault_autopilot_upgrade_status,
    VAULT_AUTOPILOT_UPGRADE_VERSION = var.vault_autopilot_upgrade_version,
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-autopilot.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
