# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The consul hosts backing the vault cluster instances"
}


resource "enos_remote_exec" "add_telemetry_to_consul" {
  for_each = var.hosts

  scripts = [abspath("${path.module}/scripts/add-consul-telemetry.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

module "restart_consul_nodes" {
  depends_on = [enos_remote_exec.add_telemetry_to_consul]
  source     = "../../restart_consul"
  hosts      = var.hosts
}
