# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.4.0"
    }
  }
}

variable "service_name" {
  type        = string
  description = "The Vault systemd service name"
  default     = "vault"
}

variable "target_hosts" {
  description = "The target machines host addresses to use for the Vault cluster"
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
}

resource "enos_remote_exec" "shutdown_multiple_nodes" {
  for_each = var.target_hosts
  inline   = ["sudo systemctl stop ${var.service_name}.service; sleep 5"]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
