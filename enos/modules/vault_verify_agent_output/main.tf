# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_agent_template_destination" {
  type        = string
  description = "The destination of the template rendered by Agent"
}

variable "vault_agent_expected_output" {
  type        = string
  description = "The output that's expected in the rendered template at vault_agent_template_destination"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

locals {
  vault_instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "verify_vault_agent_output" {
  content = templatefile("${path.module}/templates/verify-vault-agent-output.sh", {
    vault_agent_template_destination = var.vault_agent_template_destination
    vault_agent_expected_output      = var.vault_agent_expected_output
    vault_instances                  = jsonencode(local.vault_instances)
  })

  transport = {
    ssh = {
      host = local.vault_instances[0].public_ip
    }
  }
}
