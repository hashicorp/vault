# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

globals {
  backend_license_path = var.backend_license_path
  backend_tag_key      = "VaultStorage"
  build_tags = {
    "ce"               = ["ui"]
    "ent"              = ["ui", "enterprise", "ent"]
    "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
    "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
    "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
  }
  distro_packages = {
    amazon_linux = ["nc"]
    leap         = ["netcat"]
    rhel         = ["nc"]
    sles         = ["netcat"]
    ubuntu       = ["netcat"]
  }
  distro_version = {
    "amazon_linux" = var.distro_version_amazon_linux
    "leap"         = var.distro_version_leap
    "rhel"         = var.distro_version_rhel
    "sles"         = var.distro_version_sles
    "ubuntu"       = var.distro_version_ubuntu
  }
  package_manager = {
    "amazon_linux" = "yum"
    "leap"         = "zypper"
    "rhel"         = "yum"
    "sles"         = "zypper"
    "ubuntu"       = "apt"
  }
  packages = ["jq"]
  sample_attributes = {
    # NOTE(9/28/23): Temporarily use us-east-2 due to another networking in us-east-1
    # aws_region = ["us-east-1", "us-west-2"]
    aws_region                  = ["us-east-2", "us-west-2"]
    distro_version_amazon_linux = ["amzn2"]
    distro_version_leap         = ["15.4", "15.5"]
    distro_version_rhel         = ["8.8", "9.1"]
    distro_version_sles         = ["v15_sp4_standard", "v15_sp5_standard"]
    distro_version_ubuntu       = ["18.04", "20.04", "22.04"]
  }
  tags = merge({
    "Project Name" : var.project_name
    "Project" : "Enos",
    "Environment" : "ci"
  }, var.tags)
  vault_install_dir = {
    bundle  = "/opt/vault/bin"
    package = "/usr/bin"
  }
  vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
}
