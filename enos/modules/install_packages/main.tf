# Copyright (c) HashiCorp, Inc.
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
    # Note: though we generally use "amzn2" as our distro name for Amazon Linux 2, 
    # enos_host_info.hosts[each.key].distro returns "amzn", so that is what we reference here.
    "amzn"          = "yum"
    "opensuse-leap" = "zypper"
    "rhel"          = "yum"
    "sles"          = "zypper"
    "ubuntu"        = "apt"
  }
  distro_repos = {
    # Currently sles is the only distro that requires setting up repos before installing packages
    "sles" = {
      "15.5" = "https://download.opensuse.org/repositories/network:utilities/SLE_15_SP5/network:utilities.repo"
    }
    "rhel" = {
      "8.9" = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm"
      "9.3" = "https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm"
    }
  }
}

variable "packages" {
  type    = list(string)
  default = []
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts to install packages on"
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
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

# Set up repos for each distro (in order to install some packages, some distros
# require us to manually add the repo for that package first)
resource "enos_remote_exec" "distro_repo_setup" {
  for_each = var.hosts

  environment = {
    DISTRO          = enos_host_info.hosts[each.key].distro
    DISTRO_REPOS    = try(local.distro_repos[enos_host_info.hosts[each.key].distro][enos_host_info.hosts[each.key].distro_version], "__none")
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/distro-repo-setup.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "install_packages" {
  for_each   = var.hosts
  depends_on = [enos_remote_exec.distro_repo_setup]

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
