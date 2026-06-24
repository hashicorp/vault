# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    fyre = {
      source = "registry.terraform.io/hashicorp-forge/fyre"
    }
  }
}

locals {
  architectures = toset(["amd64", "ppc64le", "s390x"])

  distro_map = {
    "centos" = "centos"
    "rhel"   = "redhat"
    "ubuntu" = "ubuntu"
    "sles"   = "sles"
    "rocky"  = "rocky"
  }

  platform_map = {
    "amd64"   = "x"
    "s390x"   = "z"
    "ppc64le" = "pvm"
  }

  // Create os_versions map for each architecture
  os_versions = {
    for arch in local.architectures : arch => {
      centos  = coalesce(data.fyre_vm_os_available.current[arch].centos, [])
      redhat  = coalesce(data.fyre_vm_os_available.current[arch].redhat, [])
      rocky   = coalesce(data.fyre_vm_os_available.current[arch].rocky, [])
      rhel    = coalesce(data.fyre_vm_os_available.current[arch].redhat, [])
      sles    = coalesce(data.fyre_vm_os_available.current[arch].sles, [])
      ubuntu  = coalesce(data.fyre_vm_os_available.current[arch].ubuntu, [])
      windows = coalesce(data.fyre_vm_os_available.current[arch].windows, [])
    }
  }

  // Create a map that returns the correct distro name given the arch, distro and version.
  // This way the module can be used just like ec2_info.ami_ids. e.g.
  //   os_ids["amd64"]["rhel"]["10"]        => "RedHat 10.0"
  //   os_ids["s390x"]["sles"]["16.0"]      => "SLES 16.0"
  //   os_ids["ppc64le"]["ubuntu"]["24.04"] => "Ubuntu 24.04"
  os_ids = {
    for arch in local.architectures : arch => {
      for distro, source in local.distro_map : distro => {
        for version in distinct([
          for os_name in lookup(local.os_versions[arch], source, []) :
          trimprefix(trimsuffix(join(" ", slice(split(" ", os_name), 1, length(split(" ", os_name)))), " FIPS"), "Stream ")
        ]) :
        version => one([
          for os_name in lookup(local.os_versions[arch], source, []) : os_name
          if trimprefix(trimsuffix(join(" ", slice(split(" ", os_name), 1, length(split(" ", os_name)))), " FIPS"), "Stream ") == version
        ])
      }
    }
  }
}

data "fyre_vm_os_available" "current" {
  for_each = local.architectures
  platform = local.platform_map[each.value]
}

output "default_size" {
  description = "Default and maximum sizing metadata for each Fyre platform/site"
  value = {
    for arch in local.architectures : arch => data.fyre_vm_os_available.current[arch].default_size
  }
}

output "os_versions" {
  description = "Available OS versions grouped by architecture and distro family"
  value       = local.os_versions
}

output "os_ids" {
  description = "Available OS names grouped by architecture, normalized distro and distro_version keys"
  value       = local.os_ids
}
