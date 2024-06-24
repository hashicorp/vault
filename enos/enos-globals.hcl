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
  config_modes    = ["env", "file"]
  consul_editions = ["ce", "ent"]
  consul_versions = ["1.14.11", "1.15.7", "1.16.3", "1.17.0"]
  distros         = ["amzn2", "leap", "rhel", "sles", "ubuntu"]
  # Different distros may require different packages, or use different aliases for the same package
  distro_packages = {
    amzn2 = ["nc"]
    leap  = ["netcat", "openssl"]
    rhel  = ["nc"]
    # When installing Vault RPM packages on a SLES AMI, the openssl package provided
    # isn't named "openssl, which rpm doesn't know how to handle. Therefore we add the
    # "correctly" named one in our package installation before installing Vault.
    sles   = ["netcat-openbsd", "openssl"]
    ubuntu = ["netcat"]
  }
  distro_version = {
    "amzn2"  = var.distro_version_amzn2
    "leap"   = var.distro_version_leap
    "rhel"   = var.distro_version_rhel
    "sles"   = var.distro_version_sles
    "ubuntu" = var.distro_version_ubuntu
  }
  editions            = ["ce", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
  enterprise_editions = [for e in global.editions : e if e != "ce"]
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
    distro_version_leap   = ["15.5"]
    distro_version_rhel   = ["8.9", "9.3"]
    distro_version_sles   = ["v15_sp5_standard"]
    distro_version_ubuntu = ["20.04", "22.04"]
  }
  seals = ["awskms", "pkcs11", "shamir"]
  tags = merge({
    "Project Name" : var.project_name
    "Project" : "Enos",
    "Environment" : "ci"
  }, var.tags)
  // This reads the VERSION file, strips any pre-release metadata, and selects only initial
  // versions that are less than our current version. E.g. A VERSION file containing 1.17.0-beta2
  // would render: semverconstraint(v, "<1.17.0-0")
  upgrade_version_stripped = join("-", [split("-", chomp(file("../version/VERSION")))[0], "0"])
  // NOTE: when backporting, make sure that our initial versions are less than that
  // release branch's version. Also beware if adding versions below 1.11.x. Some scenarios
  // that use this global might not work as expected with earlier versions. Below 1.8.x is
  // not supported in any way.
  upgrade_all_initial_versions_ce  = ["1.8.12", "1.9.10", "1.10.11", "1.11.12", "1.12.11", "1.13.13", "1.14.10", "1.15.6", "1.16.3", "1.17.0"]
  upgrade_all_initial_versions_ent = ["1.8.12", "1.9.10", "1.10.11", "1.11.12", "1.12.11", "1.13.13", "1.14.13", "1.15.10", "1.16.4", "1.17.0"]
  upgrade_initial_versions_ce      = [for v in global.upgrade_all_initial_versions_ce : v if semverconstraint(v, "<${global.upgrade_version_stripped}")]
  upgrade_initial_versions_ent     = [for v in global.upgrade_all_initial_versions_ent : v if semverconstraint(v, "<${global.upgrade_version_stripped}")]
  vault_install_dir = {
    bundle  = "/opt/vault/bin"
    package = "/usr/bin"
  }
  vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
}
