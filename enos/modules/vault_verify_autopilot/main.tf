# Copyright IBM Corp. 2016, 2026
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

locals {
  leader_host = [for key, host in var.hosts : host if host.public_ip != ""][0]
}

module "verify_autopilot_status" {
  source = "../vault_run_blackbox_test"

  test_package = "github.com/hashicorp/vault/vault/external_tests/blackbox/isolated/verify"
  test_names   = ["TestAutopilotUpgradeStatus"]

  vault_addr       = null # Let it default to leader_public_ip - test runs locally, not on remote host
  vault_root_token = var.vault_root_token
  vault_edition    = "ent"

  leader_host      = local.leader_host
  leader_public_ip = local.leader_host.public_ip

  test_env_vars = {
    VAULT_AUTOPILOT_UPGRADE_STATUS  = var.vault_autopilot_upgrade_status
    VAULT_AUTOPILOT_UPGRADE_VERSION = var.vault_autopilot_upgrade_version
  }
}

module "verify_autopilot_status_output" {
  source = "../vault_run_blackbox_test"

  test_package = "github.com/hashicorp/vault/vault/external_tests/blackbox/isolated/verify"
  test_names   = ["TestAutopilotUpgradeStatusOutput"]

  vault_addr       = null # Let it default to leader_public_ip - test runs locally, not on remote host
  vault_root_token = var.vault_root_token
  vault_edition    = "ent"

  leader_host      = local.leader_host
  leader_public_ip = local.leader_host.public_ip

  test_env_vars = {
    VAULT_AUTOPILOT_UPGRADE_STATUS  = var.vault_autopilot_upgrade_status
    VAULT_AUTOPILOT_UPGRADE_VERSION = var.vault_autopilot_upgrade_version
  }
}