locals {

  // file name extensions for the install packages of vault for the various architectures, distributions and editions
  package_extensions = {
    amd64 = {
      ubuntu = {
        # "oss"     = "-1_amd64.deb"
        "ent"     = "-1_amd64.deb"
        "ent.hsm" = "-1_amd64.deb"
      }
      rhel = {
        # "oss"     = "-1.x86_64.rpm"
        "ent"     = "-1.x86_64.rpm"
        "ent.hsm" = "-1.x86_64.rpm"
      }
    }
    arm64 = {
      ubuntu = {
        # "oss" = "-1_arm64.deb"
        "ent" = "-1_arm64.deb"
      }
      rhel = {
        # "oss" = "-1.aarch64.rpm"
        "ent" = "-1.aarch64.rpm"
      }
    }
  }

  // file name prefixes for the install packages of vault for the various distributions and artifact types (package or bundle)
  artifact_package_release_names = {
    ubuntu = {
      "oss"     = "vault_"
      "ent"     = "vault-enterprise_",
      "ent.hsm" = "vault-enterprise-hsm_",
    },
    rhel = {
      "oss"     = "vault-"
      "ent"     = "vault-enterprise-",
      "ent.hsm" = "vault-enterprise-hsm-",
    }
  }

  artifact_types = ["package", "bundle"]

  // edition --> artifact name edition
  artifact_name_edition = {
    "oss"              = ""
    "ent"              = ""
    "ent.hsm"          = ".hsm"
    "ent.fips1402"     = ".fips1402"
    "ent.hsm.fips1402" = ".hsm.fips1402"
  }

  artifact_name_prefix    = var.artifact_type == "package" ? local.artifact_package_release_names[var.distro][var.edition] : "vault_"
  artifact_name_extension = var.artifact_type == "package" ? local.package_extensions[var.arch][var.distro][var.edition] : "${local.artifact_name_edition[var.edition]}_linux_${var.arch}.zip"
  artifact_name           = var.artifact_type == "package" ? "${local.artifact_name_prefix}${replace(var.product_version, "-", "~")}${local.artifact_name_extension}" : "${local.artifact_name_prefix}${var.product_version}${local.artifact_name_extension}"
}
