# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts to install packages on"
}

resource "enos_remote_exec" "make_selinux_permissive" {
  for_each = var.hosts

  scripts = [abspath("${path.module}/scripts/make-selinux-permissive.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
