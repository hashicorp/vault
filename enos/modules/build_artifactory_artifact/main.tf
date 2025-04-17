# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.2.3"
    }
  }
}

variable "artifactory_username" {
  type        = string
  description = "The username to use when connecting to artifactory"
  default     = null
}

variable "artifactory_token" {
  type        = string
  description = "The token to use when connecting to artifactory"
  default     = null
  sensitive   = true
}

variable "artifactory_host" {
  type        = string
  description = "The artifactory host to search for vault artifacts"
  default     = "https://artifactory.hashicorp.engineering/artifactory"
}

variable "artifactory_repo" {
  type        = string
  description = "The artifactory repo to search for vault artifacts"
  default     = "hashicorp-crt-stable-local*"
}

variable "arch" {}
variable "artifact_type" {}
variable "artifact_path" {}
variable "distro" {}
variable "edition" {}
variable "revision" {}
variable "product_version" {}
variable "build_tags" { default = null }
variable "bundle_path" { default = null }
variable "goarch" { default = null }
variable "goos" { default = null }

module "artifact_metadata" {
  source = "../artifact/metadata"

  arch          = var.arch
  distro        = var.distro
  edition       = var.edition
  package_type  = var.artifact_type
  vault_version = var.product_version
}

data "enos_artifactory_item" "vault" {
  username = var.artifactory_username
  token    = var.artifactory_token
  name     = module.artifact_metadata.artifact_name
  host     = var.artifactory_host
  repo     = var.artifactory_repo
  path     = "${module.artifact_metadata.product_name}/*"
  properties = tomap({
    "commit"          = var.revision,
    "product-name"    = module.artifact_metadata.product_name,
    "product-version" = module.artifact_metadata.product_version,
  })
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

output "vault_artifactory_release" {
  value = {
    url      = data.enos_artifactory_item.vault.results[0].url
    sha256   = data.enos_artifactory_item.vault.results[0].sha256
    username = var.artifactory_username
    token    = var.artifactory_token
  }
}
