# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.4.9"
    }
  }
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts for whom default softhsm configuration will be applied"
}

variable "skip" {
  type        = bool
  default     = false
  description = "Whether or not to skip initializing softhsm"
}

locals {
  // The location on disk to write the softhsm tokens to
  token_dir = "/var/lib/softhsm/tokens"

  // Where the default configuration is
  config_paths = {
    "rhel"   = "/etc/softhsm2.conf"
    "ubuntu" = "/etc/softhsm/softhsm2.conf"
  }

  host_key    = element(keys(enos_host_info.hosts), 0)
  config_path = local.config_paths[enos_host_info.hosts[local.host_key].distro]
}

resource "enos_host_info" "hosts" {
  for_each = var.hosts

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "init_softhsm" {
  for_each   = var.hosts
  depends_on = [enos_host_info.hosts]

  environment = {
    CONFIG_PATH = local.config_paths[enos_host_info.hosts[each.key].distro]
    TOKEN_DIR   = local.token_dir
    SKIP        = var.skip ? "true" : "false"
  }

  scripts = [abspath("${path.module}/scripts/init-softhsm.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

output "config_path" {
  // Technically this is actually just the first config path of our hosts.
  value = local.config_path
}

output "token_dir" {
  value = local.token_dir
}

output "skipped" {
  value = var.skip
}
