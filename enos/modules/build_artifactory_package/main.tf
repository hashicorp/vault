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
  default     = "9"
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

locals {
  // File name prefixes for the various distributions and editions
  artifact_prefix = {
    ubuntu = {
      "ce"               = "vault_"
      "ent"              = "vault-enterprise_",
      "ent.hsm"          = "vault-enterprise-hsm_",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402_",
      "oss"              = "vault_"
    },
    rhel = {
      "ce"               = "vault-"
      "ent"              = "vault-enterprise-",
      "ent.hsm"          = "vault-enterprise-hsm-",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
      "oss"              = "vault-"
    }
  }

  // Format the version and edition to use in the artifact name
  artifact_version = {
    "ce"               = "${var.product_version}"
    "ent"              = "${var.product_version}+ent"
    "ent.hsm"          = "${var.product_version}+ent"
    "ent.hsm.fips1402" = "${var.product_version}+ent"
    "oss"              = "${var.product_version}"
  }

  // File name extensions for the various architectures and distributions
  artifact_extension = {
    amd64 = {
      ubuntu = "-1_amd64.deb"
      rhel   = "-1.x86_64.rpm"
    }
    arm64 = {
      ubuntu = "-1_arm64.deb"
      rhel   = "-1.aarch64.rpm"
    }
  }

  // Use the above variables to construct the artifact name to look up in Artifactory.
  // Will look something like:
  //  vault_1.12.2-1_arm64.deb
  //  vault-enterprise_1.12.2+ent-1_amd64.deb
  //  vault-enterprise-hsm-1.12.2+ent-1.x86_64.rpm
  artifact_name = "${local.artifact_prefix[var.distro][var.edition]}${local.artifact_version[var.edition]}${local.artifact_extension[var.arch][var.distro]}"

  // The path within the Artifactory repo that corresponds to the appropriate architecture
  artifactory_repo_path_dir = {
    "amd64" = "x86_64"
    "arm64" = "aarch64"
  }
}

data "enos_artifactory_item" "vault_package" {
  username = var.artifactory_username
  token    = var.artifactory_token
  name     = local.artifact_name
  host     = var.artifactory_host
  repo     = var.distro == "rhel" ? "hashicorp-rpm-release-local*" : "hashicorp-apt-release-local*"
  path     = var.distro == "rhel" ? "RHEL/${var.distro_version}/${local.artifactory_repo_path_dir[var.arch]}/stable" : "pool/${var.arch}/main"
}

output "results" {
  value = data.enos_artifactory_item.vault_package.results
}

output "url" {
  value       = data.enos_artifactory_item.vault_package.results[0].url
  description = "The artifactory download url for the artifact"
}

output "sha256" {
  value       = data.enos_artifactory_item.vault_package.results[0].sha256
  description = "The sha256 checksum for the artifact"
}

output "size" {
  value       = data.enos_artifactory_item.vault_package.results[0].size
  description = "The size in bytes of the artifact"
}

output "name" {
  value       = data.enos_artifactory_item.vault_package.results[0].name
  description = "The name of the artifact"
}

output "release" {
  value = {
    url      = data.enos_artifactory_item.vault_package.results[0].url
    sha256   = data.enos_artifactory_item.vault_package.results[0].sha256
    username = var.artifactory_username
    token    = var.artifactory_token
  }
}
