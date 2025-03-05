# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

// Given the architecture, distro, version, edition, and desired package type,
// return the metadata for an artifact.

variable "arch" {
  description = "The platform architecture"
  type        = string

  validation {
    condition     = contains(["amd64", "arm64", "s390x"], var.arch)
    error_message = <<-EOF
      distro must be one of "amd64", "arm64", "s390x"
    EOF
  }
}

variable "distro" {
  description = "The operating system distro"
  type        = string

  validation {
    condition     = contains(["amzn", "leap", "rhel", "sles", "ubuntu"], var.distro)
    error_message = <<-EOF
      distro must be one of "amzn", "leap", "rhel", "sles", "ubuntu"
    EOF
  }
}

variable "edition" {
  description = "The Vault edition. E.g. ent or ent.hsm.fips1403"
  type        = string

  validation {
    condition     = contains(["oss", "ce", "ent", "ent.fips1402", "ent.fips1403", "ent.hsm", "ent.hsm.fips1402", "ent.hsm.fips1403"], var.edition)
    error_message = <<-EOF
      edition must be one of "oss", "ce", "ent", "ent.fips1402", "ent.fips1403", "ent.hsm", "ent.hsm.fips1402", "ent.hsm.fips1403"
    EOF
  }
}

variable "package_type" {
  description = "The preferred packaging"
  type        = string

  validation {
    condition     = contains(["package", "rpm", "deb", "zip", "bundle"], var.package_type)
    error_message = <<-EOF
      package_type must be one of "package", "rpm", "deb", "zip", "bundle"
    EOF
  }
}

variable "vault_version" {
  description = "The of Vault or Vault Enterprise. E.g 1.18.2, 1.19.0-rc1, 1.18.5+ent.hsm"
  type        = string
}

locals {
  package_extension_amd64_deb = "-1_amd64.deb"
  package_extension_amd64_rpm = "-1.x86_64.rpm"
  package_extension_arm64_deb = "-1_arm64.deb"
  package_extension_arm64_rpm = "-1.aarch64.rpm"
  package_extension_s390x_deb = "-1_s390x.deb"
  package_extension_s390x_rpm = "-1.s390x.rpm"

  // file name extensions for the install packages of vault for the various architectures, distributions and editions
  package_extensions = {
    amd64 = {
      amzn   = local.package_extension_amd64_rpm
      leap   = local.package_extension_amd64_rpm
      rhel   = local.package_extension_amd64_rpm
      sles   = local.package_extension_amd64_rpm
      ubuntu = local.package_extension_amd64_deb
    }
    arm64 = {
      amzn   = local.package_extension_arm64_rpm
      leap   = local.package_extension_arm64_rpm
      rhel   = local.package_extension_arm64_rpm
      sles   = local.package_extension_arm64_rpm
      ubuntu = local.package_extension_arm64_deb
    }
    s390x = {
      amzn   = null
      leap   = local.package_extension_s390x_rpm
      rhel   = local.package_extension_s390x_rpm
      sles   = local.package_extension_s390x_rpm
      ubuntu = local.package_extension_s390x_deb
    }
  }

  package_prefixes_rpm = {
    "ce"               = "vault-"
    "ent"              = "vault-enterprise-",
    "ent.fips1402"     = "vault-enterprise-fips1402-",
    "ent.fips1403"     = "vault-enterprise-fips1403-",
    "ent.hsm"          = "vault-enterprise-hsm-",
    "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
    "ent.hsm.fips1403" = "vault-enterprise-hsm-fips1403-",
    "oss"              = "vault-"
  }

  package_prefixes_deb = {
    "ce"               = "vault_"
    "ent"              = "vault-enterprise_",
    "ent.fips1402"     = "vault-enterprise-fips1402_",
    "ent.fips1403"     = "vault-enterprise-fips1403_",
    "ent.hsm"          = "vault-enterprise-hsm_",
    "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402_",
    "ent.hsm.fips1403" = "vault-enterprise-hsm-fips1403_",
    "oss"              = "vault_"
  }

  // file name prefixes for the install packages of vault for the various distributions and artifact types (package or bundle)
  package_prefixes = {
    amzn   = local.package_prefixes_rpm,
    leap   = local.package_prefixes_rpm,
    rhel   = local.package_prefixes_rpm,
    sles   = local.package_prefixes_rpm,
    ubuntu = local.package_prefixes_deb,
  }

  // Stable release Artifactory repos for packages
  release_repo_rpm = "hashicorp-rpm-release-local*"
  release_repo_apt = "hashicorp-apt-release-local*"
  release_repos = {
    amzn   = local.release_repo_rpm
    rhel   = local.release_repo_rpm
    ubuntu = local.release_repo_apt
  }
  release_repo = local.release_repos[var.distro]

  // Stable release Artifactory paths for packages
  release_package_rpm_arch = {
    "amd64" = "x86_64",
    "arm64" = "aarch64",
    "s390x" = "s390x",
  }
  release_path_deb     = "pool/${var.arch}/main"
  release_sub_path_rpm = "${local.release_package_rpm_arch[var.arch]}/stable"
  release_path_distro = {
    amzn = {
      "2"      = "AmazonLinux/2/${local.release_sub_path_rpm}"
      "2023"   = "AmazonLinux/latest/${local.release_sub_path_rpm}"
      "latest" = "AmazonLinux/latest/${local.release_sub_path_rpm}"
    }
    leap = {
      "15.6" = "RHEL/9/${local.release_sub_path_rpm}"
    }
    rhel = {
      "8" = "RHEL/8/${local.release_sub_path_rpm}"
      "9" = "RHEL/9/${local.release_sub_path_rpm}"
    }
    sles = {
      "15.6" = "RHEL/9/${local.release_sub_path_rpm}"
    }
    ubuntu = {
      "20.04" = local.release_path_deb,
      "22.04" = local.release_path_deb,
      "24.04" = local.release_path_deb,
    }
  }
  release_paths = local.release_path_distro[var.distro]

  // Reduce our supported inputs into two classes: system packages or a binary bundled into a zip archive.
  package_type = contains(["package", "deb", "rpm"], var.package_type) ? "package" : "bundle"

  // Get the base version. This might still include pre-release metadata
  // E.g. 1.18.2 => 1.18.2, 1.18.0-rc1 => 1.18.0-rc1, 1.18.0+ent.hsm => 1.18.0
  semverish_version = try(split("+", var.vault_version)[0], var.vault_version)

  // Determine the "product name". This corresponds properties on the artifactory artifact.
  product_name = strcontains(var.edition, "ent") ? "vault-enterprise" : "vault"

  // Create the "product version", which is corresponds to properties on the artifactory artifact.
  // It's the version along with edition metadata. We normalize all enterprise editions to .ent.
  // E.g. 1.16.0-beta1+ent.hsm.fips1403 -> 1.16.0-beta+ent
  product_version = strcontains(var.edition, "ent") ? "${local.semverish_version}+ent" : local.semverish_version

  // Convert product version strings to a syntax that matches deb and rpm packaging.
  // E.g. 1.16.0-beta+ent -> 1.16.0~beta+ent
  package_version = replace(local.product_version, "-", "~")

  // Get the bundle version. If the vault_version includes metadata, use it. Otherwise add the edition to it.
  bundle_version = strcontains(var.vault_version, "+") ? var.vault_version : strcontains(var.edition, "ent") ? "${var.vault_version}+${var.edition}" : var.vault_version

  // Prefix for the artifact name. E.g.: vault_, vault-, vault-enterprise_, vault-enterprise-hsm-fips1402-, etc
  artifact_name_prefix = local.package_type == "package" ? local.package_prefixes[var.distro][var.edition] : "vault_"

  // The version for the artifact name.
  artifact_version = local.package_type == "package" ? local.package_version : local.bundle_version

  // Suffix and extension for the artifact name. E.g.: _linux_<arch>.zip,
  artifact_name_extension = local.package_type == "package" ? local.package_extensions[var.arch][var.distro] : "_linux_${var.arch}.zip"

  // Combine prefix/suffix/extension together to form the artifact name
  artifact_name = "${local.artifact_name_prefix}${local.artifact_version}${local.artifact_name_extension}"

}
output "artifact_name" {
  value = local.artifact_name
}

output "package_type" {
  value = local.package_type
}

output "package_version" {
  value = local.package_version
}

output "product_name" {
  value = local.product_name
}

output "product_version" {
  value = local.product_version
}

output "release_repo" {
  value = local.release_repo
}

output "release_paths" {
  value = local.release_paths
}
