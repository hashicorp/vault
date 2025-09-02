# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The node to shut down"
}

resource "enos_remote_exec" "shutdown_node" {
  inline = ["sudo shutdown -P --no-wall; exit 0"]

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}
