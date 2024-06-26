# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {

  // file name extensions for the install packages of vault for the various architectures, distributions and editions
  package_extensions = {
    amd64 = {
      amzn2  = "-1.x86_64.rpm"
      leap   = "-1.x86_64.rpm"
      rhel   = "-1.x86_64.rpm"
      sles   = "-1.x86_64.rpm"
      ubuntu = "-1_amd64.deb"
    }
    arm64 = {
      amzn2  = "-1.aarch64.rpm"
      leap   = "-1.aarch64.rpm"
      rhel   = "-1.aarch64.rpm"
      sles   = "-1.aarch64.rpm"
      ubuntu = "-1_arm64.deb"
    }
  }

  // product_version --> artifact_version
  artifact_version = replace(var.product_version, var.edition, "ent")

  // file name prefixes for the install packages of vault for the various distributions and artifact types (package or bundle)
  artifact_package_release_names = {
    amzn2 = {
      "ce"               = "vault-"
      "ent"              = "vault-enterprise-",
      "ent.fips1402"     = "vault-enterprise-fips1402-",
      "ent.hsm"          = "vault-enterprise-hsm-",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
    },
    leap = {
      "ce"               = "vault-"
      "ent"              = "vault-enterprise-",
      "ent.fips1402"     = "vault-enterprise-fips1402-",
      "ent.hsm"          = "vault-enterprise-hsm-",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
    },
    rhel = {
      "ce"               = "vault-"
      "ent"              = "vault-enterprise-",
      "ent.fips1402"     = "vault-enterprise-fips1402-",
      "ent.hsm"          = "vault-enterprise-hsm-",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
    },
    sles = {
      "ce"               = "vault-"
      "ent"              = "vault-enterprise-",
      "ent.fips1402"     = "vault-enterprise-fips1402-",
      "ent.hsm"          = "vault-enterprise-hsm-",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402-",
    }
    ubuntu = {
      "ce"               = "vault_"
      "ent"              = "vault-enterprise_",
      "ent.fips1402"     = "vault-enterprise-fips1402_",
      "ent.hsm"          = "vault-enterprise-hsm_",
      "ent.hsm.fips1402" = "vault-enterprise-hsm-fips1402_",
    }
  }

  # Prefix for the artifact name. Ex: vault_, vault-, vault-enterprise_, vault-enterprise-hsm-fips1402-, etc
  artifact_name_prefix = var.artifact_type == "package" ? local.artifact_package_release_names[var.distro][var.edition] : "vault_"
  # Suffix and extension for the artifact name. Ex: _linux_<arch>.zip, 
  artifact_name_extension = var.artifact_type == "package" ? local.package_extensions[var.arch][var.distro] : "_linux_${var.arch}.zip"
  # Combine prefix/suffix/extension together to form the artifact name
  artifact_name = var.artifact_type == "package" ? "${local.artifact_name_prefix}${replace(local.artifact_version, "-", "~")}${local.artifact_name_extension}" : "${local.artifact_name_prefix}${var.product_version}${local.artifact_name_extension}"
}
