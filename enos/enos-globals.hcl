# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

globals {
  archs            = ["amd64", "arm64"]
  artifact_sources = ["local", "crt", "artifactory"]
  artifact_types   = ["bundle", "package"]
  backends         = ["consul", "raft"]
  backend_tag_key  = "VaultStorage"
  build_tags = {
    "ce"               = ["ui"]
    "ent"              = ["ui", "enterprise", "ent"]
    "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
    "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
    "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
  }
  config_modes    = ["env", "file"]
  consul_versions = ["1.14.11", "1.15.7", "1.16.3", "1.17.0"]
  distros         = ["ubuntu", "rhel"]
  distro_version = {
    "rhel"   = var.rhel_distro_version
    "ubuntu" = var.ubuntu_distro_version
  }
  editions = ["ce", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
  packages = ["jq"]
  distro_packages = {
    ubuntu = ["netcat"]
    rhel   = ["nc"]
  }
  sample_attributes = {
    aws_region = ["us-east-1", "us-west-2"]
  }
  seals = ["awskms", "pkcs11", "shamir"]
  tags = merge({
    "Project Name" : var.project_name
    "Project" : "Enos",
    "Environment" : "ci"
  }, var.tags)
  // NOTE: when backporting, make sure that our initial versions are less than that
  // release branch's version. Also beware if adding versions below 1.11.x. Some scenarios
  // that use this global might not work as expected with earlier versions. Below 1.8.x is
  // not supported in any way.
  upgrade_initial_versions = ["1.11.12", "1.12.11", "1.13.11", "1.14.7", "1.15.3"]
  vault_install_dir_packages = {
    rhel   = "/bin"
    ubuntu = "/usr/bin"
  }
  vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
}
