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
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_agent_expected_output" {
  type        = string
  description = "The output that's expected in the rendered template at vault_agent_template_destination"
}

variable "vault_agent_template_destination" {
  type        = string
  description = "The destination of the template rendered by Agent"
}

resource "enos_remote_exec" "verify_vault_agent_output" {
  environment = {
    VAULT_AGENT_TEMPLATE_DESTINATION = var.vault_agent_template_destination
    VAULT_AGENT_EXPECTED_OUTPUT      = var.vault_agent_expected_output
  }

  scripts = [abspath("${path.module}/scripts/verify-vault-agent-output.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}
