# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "primary_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The primary cluster leader host"
}

variable "secondary_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The secondary cluster leader host"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "primary_root_token" {
  type        = string
  description = "The root token for the primary cluster"
  default     = ""
}

variable "secondary_root_token" {
  type        = string
  description = "The root token for the secondary cluster"
  default     = ""
}

variable "wrapping_token" {
  type        = string
  description = "The wrapping token created on primary cluster"
  default     = null
}

locals {
  primary_leader_addr   = var.ip_version == 6 ? var.primary_leader_host.ipv6 : var.primary_leader_host.private_ip
  secondary_leader_addr = var.ip_version == 6 ? var.secondary_leader_host.ipv6 : var.secondary_leader_host.private_ip

  # Extract JSON output from test stdout
  # The test prints JSON to stdout, which we need to extract and parse
  primary_json_output = try(
    regex("(?s)\\{.*\"mode\".*\\}", module.verify_replication_status_on_primary.test_result)[0],
    "{}"
  )
  secondary_json_output = try(
    regex("(?s)\\{.*\"mode\".*\\}", module.verify_replication_status_on_secondary.test_result)[0],
    "{}"
  )

  primary_replication_status   = jsondecode(local.primary_json_output)
  secondary_replication_status = jsondecode(local.secondary_json_output)
}

module "verify_replication_status_on_primary" {
  source = "../vault_run_blackbox_test"

  leader_host       = var.primary_leader_host
  leader_public_ip  = var.primary_leader_host.public_ip
  vault_root_token  = var.primary_root_token
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
  test_package      = "./vault/external_tests/blackbox/isolated/verify"
  test_names        = ["TestDRReplicationStatusOutput"]
  vault_edition     = "ent" # DR replication is enterprise-only

  test_env_vars = {
    PRIMARY_LEADER_ADDR = local.primary_leader_addr
  }
}

module "verify_replication_status_on_secondary" {
  source = "../vault_run_blackbox_test"

  leader_host       = var.secondary_leader_host
  leader_public_ip  = var.secondary_leader_host.public_ip
  vault_root_token  = var.secondary_root_token
  vault_addr        = var.vault_addr
  vault_install_dir = var.vault_install_dir
  test_package      = "./vault/external_tests/blackbox/isolated/verify"
  test_names        = ["TestDRReplicationStatusOutput"]
  vault_edition     = "ent" # DR replication is enterprise-only

  test_env_vars = {
    PRIMARY_LEADER_ADDR = local.primary_leader_addr
  }
}

output "primary_replication_status" {
  value = local.primary_replication_status
}

output "known_primary_cluster_addrs" {
  value = try(local.secondary_replication_status.known_primary_cluster_addrs, [])
}

output "secondary_replication_status" {
  value = local.secondary_replication_status
}

output "primary_replication_data_secondaries" {
  value = try(local.primary_replication_status.secondaries, [])
}

output "secondary_replication_data_primaries" {
  value = try(local.secondary_replication_status.primaries, [])
}