# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "arch" {
  type        = string
  description = "The architecture for the desired artifact"
}

variable "artifactory_username" {
  type        = string
  description = "The username to use when connecting to Artifactory"
}

variable "artifactory_token" {
  type        = string
  description = "The token to use when connecting to Artifactory"
  sensitive   = true
}

variable "artifactory_host" {
  type        = string
  description = "The Artifactory host to search for Vault artifacts"
  default     = "https://artifactory.hashicorp.engineering/artifactory"
}

variable "distro" {
  type        = string
  description = "The distro for the desired artifact (ubuntu or rhel)"
}

variable "distro_version" {
  type        = string
  description = "The RHEL version for .rpm packages"
}

variable "edition" {
  type        = string
  description = "The edition of Vault to use"
}

variable "product_version" {
  type        = string
  description = "The version of Vault to use"
}

// Shim variables that we don't use but include to satisfy the build module "interface"
variable "artifact_path" { default = null }
variable "artifact_type" { default = null }
variable "artifactory_repo" { default = null }
variable "build_tags" { default = null }
variable "build_ui" { default = null }
variable "bundle_path" { default = null }
variable "goarch" { default = null }
variable "goos" { default = null }
variable "revision" { default = null }

module "artifact_metadata" {
  source = "../artifact/metadata"

  arch          = var.arch
  distro        = var.distro
  edition       = var.edition
  package_type  = var.artifact_type != null ? var.artifact_type : "package"
  vault_version = var.product_version
}

data "enos_artifactory_item" "vault" {
  username = var.artifactory_username
  token    = var.artifactory_token
  name     = module.artifact_metadata.artifact_name
  host     = var.artifactory_host
  repo     = module.artifact_metadata.release_repo
  path     = module.artifact_metadata.release_paths[var.distro_version]
}

output "results" {
  value = data.enos_artifactory_item.vault.results
}

output "url" {
  value       = data.enos_artifactory_item.vault.results[0].url
  description = "The artifactory download url for the artifact"
}

output "sha256" {
  value       = data.enos_artifactory_item.vault.results[0].sha256
  description = "The sha256 checksum for the artifact"
}

output "size" {
  value       = data.enos_artifactory_item.vault.results[0].size
  description = "The size in bytes of the artifact"
}

output "name" {
  value       = data.enos_artifactory_item.vault.results[0].name
  description = "The name of the artifact"
}

output "release" {
  value = {
    url      = data.enos_artifactory_item.vault.results[0].url
    sha256   = data.enos_artifactory_item.vault.results[0].sha256
    username = var.artifactory_username
    token    = var.artifactory_token
  }
}
