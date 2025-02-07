# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "audit_log_file_path" {
  type = string
}

variable "leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The cluster leader host. Only the leader write to the audit log"
}

variable "radar_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = "/opt/vault-radar/bin"
}

variable "radar_license_path" {
  description = "The path to a vault-radar license file"
}

variable "radar_version" {
  description = "The version of Vault Radar to install"
  default     = "0.17.0" # must be >= 0.17.0
  // NOTE: A `semverconstraint` validation condition would be very useful here
  // when we get around to exporting our custom enos funcs in the provider.
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "vault_unit_name" {
  type        = string
  description = "The vault unit name"
  default     = "vault"
}

resource "enos_bundle_install" "radar" {
  destination = var.radar_install_dir

  release = {
    product = "vault-radar"
    version = var.radar_version
    // Radar doesn't have CE/Ent editions. CE is equivalent to no edition metadata.
    edition = "ce"
  }

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

resource "enos_remote_exec" "scan_logs_for_secrets" {
  depends_on = [
    enos_bundle_install.radar,
  ]

  environment = {
    AUDIT_LOG_FILE_PATH     = var.audit_log_file_path
    VAULT_ADDR              = var.vault_addr
    VAULT_RADAR_INSTALL_DIR = var.radar_install_dir
    VAULT_RADAR_LICENSE     = file(var.radar_license_path)
    VAULT_TOKEN             = var.vault_root_token
    VAULT_UNIT_NAME         = var.vault_unit_name
  }

  scripts = [abspath("${path.module}/scripts/scan_logs_for_secrets.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
