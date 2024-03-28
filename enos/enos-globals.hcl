# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

globals {
  archs                = ["amd64", "arm64"]
  artifact_sources     = ["local", "crt", "artifactory"]
  artifact_types       = ["bundle", "package"]
  backends             = ["consul", "raft"]
  backend_license_path = abspath(var.backend_license_path != null ? var.backend_license_path : joinpath(path.root, "./support/consul.hclic"))
  backend_tag_key      = "VaultStorage"
  build_tags = {
    "ce"               = ["ui"]
    "ent"              = ["ui", "enterprise", "ent"]
    "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
    "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
    "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
  }
  consul_editions = ["ce", "ent"]
  consul_versions = ["1.14.11", "1.15.7", "1.16.3", "1.17.0"]
  distros         = ["amzn2", "leap", "rhel", "sles", "ubuntu"]
  # Different distros may require different packages, or use different aliases for the same package
  distro_packages = {
    amzn2 = ["nc"]
    leap  = ["netcat", "openssl"]
    rhel  = ["nc"]
    # When installing Vault RPM packages, SLES searches for openssl by a different name
    # than the one that comes pre-installed on the AMI. Therefore we add the
    # "correctly" named one in our package installation before installing Vault.
    sles   = ["nc", "openssl"]
    ubuntu = ["netcat"]
  }
  distro_version = {
    "amzn2"  = var.distro_version_amzn2
    "leap"   = var.distro_version_leap
    "rhel"   = var.distro_version_rhel
    "sles"   = var.distro_version_sles
    "ubuntu" = var.distro_version_ubuntu
  }
  editions = ["ce", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
  package_manager = {
    "amzn2"  = "yum"
    "leap"   = "zypper"
    "rhel"   = "yum"
    "sles"   = "zypper"
    "ubuntu" = "apt"
  }
  packages = ["jq"]
  sample_attributes = {
    aws_region            = ["us-east-1", "us-west-2"]
    distro_version_amzn2  = ["2"]
    distro_version_leap   = ["15.4", "15.5"]
    distro_version_rhel   = ["8.8", "9.1"]
    distro_version_sles   = ["v15_sp5_standard"]
    distro_version_ubuntu = ["18.04", "20.04", "22.04"]
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
  vault_install_dir = {
    bundle  = "/opt/vault/bin"
    package = "/usr/bin"
  }
  vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
}
