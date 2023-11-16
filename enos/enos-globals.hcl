# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

globals {
  backend_tag_key = "VaultStorage"
  build_tags = {
    "ce"               = ["ui"]
    "ent"              = ["ui", "enterprise", "ent"]
    "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
    "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
    "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
  }
  distro_version = {
    "rhel"   = var.rhel_distro_version
    "ubuntu" = var.ubuntu_distro_version
  }
  packages = ["jq"]
  distro_packages = {
    ubuntu = ["netcat"]
    rhel   = ["nc"]
  }
  sample_attributes = {
    # NOTE(9/28/23): Temporarily use us-east-2 due to another networking in us-east-1
    # aws_region = ["us-east-1", "us-west-2"]
    aws_region = ["us-east-2", "us-west-2"]
  }
  tags = merge({
    "Project Name" : var.project_name
    "Project" : "Enos",
    "Environment" : "ci"
  }, var.tags)
  vault_install_dir_packages = {
    rhel   = "/bin"
    ubuntu = "/usr/bin"
  }
  vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
}
