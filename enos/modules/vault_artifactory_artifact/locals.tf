# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {

  // file name extensions for the install packages of vault for the various architectures, distributions and editions
  package_extensions = {
    amd64 = {
      ubuntu = "-1_amd64.deb"
      rhel   = "-1.x86_64.rpm"
    }
    arm64 = {
      ubuntu = "-1_arm64.deb"
      rhel   = "-1.aarch64.rpm"
    }
  }

  // product_version --> artifact_version
  artifact_version = replace(var.product_version, var.edition, "ent")

  // file name prefixes for the install packages of vault for the various distributions and artifact types (package or bundle)
  artifact_package_release_names = {
    ubuntu = {
      "oss"              = "vault_"
      "ent"              = "vault-enterprise_",
      "ent.fips1402"     = "vault-enterprise-fips1402_",
      "ent.hsm"          = "vault-enterprise-hsm_",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402_",
    },
    rhel = {
      "oss"              = "vault-"
      "ent"              = "vault-enterprise-",
      "ent.fips1402"     = "vault-enterprise-fips1402-",
      "ent.hsm"          = "vault-enterprise-hsm-",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
    }
  }

  // edition --> artifact name edition
  artifact_name_edition = {
    "oss"              = ""
    "ent"              = ""
    "ent.hsm"          = ".hsm"
    "ent.fips1402"     = ".fips1402"
    "ent.hsm.fips1402" = ".hsm.fips1402"
  }

  artifact_name_prefix    = var.artifact_type == "package" ? local.artifact_package_release_names[var.distro][var.edition] : "vault_"
  artifact_name_extension = var.artifact_type == "package" ? local.package_extensions[var.arch][var.distro] : "_linux_${var.arch}.zip"
  artifact_name           = var.artifact_type == "package" ? "${local.artifact_name_prefix}${replace(local.artifact_version, "-", "~")}${local.artifact_name_extension}" : "${local.artifact_name_prefix}${var.product_version}${local.artifact_name_extension}"
}
