# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "cloud_provider" {
  description = "The cloud provider of the HCP HVN and Vault cluster."
  type        = string
  default     = "aws"
}

variable "cloud_region" {
  description = "The region of the HCP HVN and Vault cluster."
  type        = string
  default     = "us-west-2"
}

variable "maintenance_window_day" {
  description = "The maintenance window day"
  type        = string
  default     = "FRIDAY"
}

variable "maintenance_window_time" {
  description = "The maintenance window time"
  type        = string
  default     = "WINDOW_12PM_4PM"
}

variable "min_vault_version" {
  description = "The minimum vault version. This also corresponds to the image id"
  type        = string
  default     = null
}

variable "tier" {
  description = "Tier of the HCP Vault cluster. Valid options for tiers."
  type        = string
  // NOTE: we can't use dev for custom images
  default = "plus_small"
}

variable "upgrade_type" {
  description = "The upgrade strategy"
  type        = string
  default     = "MANUAL"
}

data "enos_environment" "localhost" {}

resource "random_string" "id" {
  length  = 4
  lower   = true
  upper   = false
  numeric = false
  special = false
}

locals {
  // Generate a unique identifier for our scenario. If we've been given a
  // min_vault_version we'll use that as it will likely be the version and
  // a SHA of a custom image. Make sure it doesn't have special characters.
  // Otherwise, just use a random string.
  id = var.min_vault_version != null ? try(replace(var.min_vault_version, "/[^0-9A-Za-z]/", "-"), random_string.id.result) : random_string.id.result
}

resource "hcp_hvn" "default" {
  hvn_id         = local.id
  cloud_provider = var.cloud_provider
  region         = var.cloud_region
}

resource "hcp_vault_cluster" "enos" {
  depends_on = [
    hcp_hvn.default,
  ]

  hvn_id            = local.id
  cluster_id        = "enos-${local.id}"
  tier              = var.tier
  public_endpoint   = true
  min_vault_version = var.min_vault_version

  dynamic "ip_allowlist" {
    for_each = data.enos_environment.localhost.public_ipv4_addresses
    content {
      address = "${ip_allowlist.value}/32"
    }
  }

  /*
  major_version_upgrade_config {
    maintenance_window_day  = var.maintenance_window_day
    maintenance_window_time = var.maintenance_window_time
    upgrade_type            = var.upgrade_type
  }
  */
}

output "cloud_provider" {
  value = hcp_vault_cluster.enos.cloud_provider
}

output "cluster_id" {
  value = hcp_vault_cluster.enos.cluster_id
}

output "created_at" {
  value = hcp_vault_cluster.enos.created_at
}

output "id" {
  value = hcp_vault_cluster.enos.id
}

output "namespace" {
  value = hcp_vault_cluster.enos.namespace
}

output "organization_id" {
  value = hcp_vault_cluster.enos.organization_id
}

output "region" {
  value = hcp_vault_cluster.enos.region
}

output "self_link" {
  value = hcp_vault_cluster.enos.self_link
}

output "state" {
  value = hcp_vault_cluster.enos.state
}

output "vault_private_endpoint_url" {
  value = hcp_vault_cluster.enos.vault_private_endpoint_url
}

output "vault_proxy_endpoint_url" {
  value = hcp_vault_cluster.enos.vault_proxy_endpoint_url
}

output "vault_public_endpoint_url" {
  value = hcp_vault_cluster.enos.vault_public_endpoint_url
}

output "vault_version" {
  value = hcp_vault_cluster.enos.vault_version
}
