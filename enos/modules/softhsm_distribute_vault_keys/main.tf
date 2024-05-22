# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.4.9"
    }
  }
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts for whom we'll distribute the softhsm tokens and keys"
}

variable "token_base64" {
  type        = string
  description = "The base64 encoded gzipped tarball of the softhsm token"
}

locals {
  // The user/group name for softhsm
  softhsm_groups = {
    "rhel"   = "ods"
    "ubuntu" = "softhsm"
  }

  // Determine if we should skip distribution. If we haven't been passed in a base64 token tarball
  // we should short circuit the rest of the module.
  skip = var.token_base64 == null || var.token_base64 == "" ? true : false
}

module "install" {
  // TODO: Should packages take a string instead of array so we can plan with unknown values that could change?
  source = "../softhsm_install"

  hosts         = var.hosts
  include_tools = false # we don't need opensc on machines that did not create the HSM.
}

module "initialize" {
  source     = "../softhsm_init"
  depends_on = [module.install]

  hosts = var.hosts
  skip  = local.skip
}

# In order for the vault service to access our keys we need to deal with ownership of files. Make
# sure we have a vault user on the machine if it doesn't already exist. Our distribution script
# below will handle adding vault to the "softhsm" group and setting ownership of the tokens.
resource "enos_user" "vault" {
  for_each = var.hosts

  name     = "vault"
  home_dir = "/etc/vault.d"
  shell    = "/bin/false"

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

// Get the host information so we can ensure that the correct user/group is used for softhsm.
resource "enos_host_info" "hosts" {
  for_each = var.hosts

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

// Distribute our softhsm token and keys to the given hosts.
resource "enos_remote_exec" "distribute_token" {
  for_each = var.hosts
  depends_on = [
    module.initialize,
    enos_user.vault,
    enos_host_info.hosts,
  ]

  environment = {
    TOKEN_BASE64  = var.token_base64
    TOKEN_DIR     = module.initialize.token_dir
    SOFTHSM_GROUP = local.softhsm_groups[enos_host_info.hosts[each.key].distro]
  }

  scripts = [abspath("${path.module}/scripts/distribute-token.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

output "lib" {
  value = module.install.lib
}
