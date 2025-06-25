# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The consul hosts"
}

resource "enos_remote_exec" "restart" {
  for_each = var.hosts

  scripts = [abspath("${path.module}/scripts/restart-consul.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

