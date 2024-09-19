# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "old_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances to be shutdown"
}

resource "enos_remote_exec" "shutdown_multiple_nodes" {
  for_each = var.old_hosts
  inline   = ["sudo shutdown -P --no-wall; exit 0"]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
