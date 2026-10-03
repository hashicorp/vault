# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

locals {
  arch = {
    "amd64" = "x86_64"
    "arm64" = "aarch64"
  }
  package_manager = {
    "amzn"   = "yum"
    "rhel"   = "dnf"
    "sles"   = "zypper"
    "ubuntu" = "apt"
  }
  distro_repos = {
    // NOTE: The versions here always correspond to the output of enos_host_info.distro_version. These are used in
    // several modules so if you change the keys here also consider the "artifact/metadata", "ec2_info",
    "sles" = {
      "15.7" = "https://download.opensuse.org/repositories/network:utilities/15.6/network:utilities.repo"
      "16.0" = "https://download.opensuse.org/repositories/network:utilities/16.0/network:utilities.repo"
    }
    "rhel" = {
      "8.10" = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm"
      "9.7"  = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm"
      "10.1" = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-10.noarch.rpm"
    }
  }
}

variable "packages" {
  type    = list(string)
  default = []
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts to install packages on"
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out. This is applied to each step so total timeout will be longer."
  default     = 120
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

resource "enos_host_info" "hosts" {
  for_each = var.hosts

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Synchronize repositories on remote machines. This does not update packages but only ensures that
# the remote hosts are configured with default upstream repositories that have been refreshed to
# the latest metedata.
resource "enos_remote_exec" "synchronize_repos" {
  for_each = var.hosts

  environment = {
    DISTRO          = enos_host_info.hosts[each.key].distro
    PACKAGE_MANAGER = local.package_manager[enos_host_info.hosts[each.key].distro]
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/synchronize-repos.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Add any additional repositories.
resource "enos_remote_exec" "add_repos" {
  for_each   = var.hosts
  depends_on = [enos_remote_exec.synchronize_repos]

  environment = {
    DISTRO_REPOS    = try(local.distro_repos[enos_host_info.hosts[each.key].distro][enos_host_info.hosts[each.key].distro_version], "__none")
    PACKAGE_MANAGER = local.package_manager[enos_host_info.hosts[each.key].distro]
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/add-repos.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Install any required packages.
resource "enos_remote_exec" "install_packages" {
  for_each = var.hosts
  depends_on = [
    enos_remote_exec.synchronize_repos,
    enos_remote_exec.add_repos,
  ]

  environment = {
    PACKAGE_MANAGER = local.package_manager[enos_host_info.hosts[each.key].distro]
    PACKAGES        = length(var.packages) >= 1 ? join(" ", var.packages) : "__skip"
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-packages.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
