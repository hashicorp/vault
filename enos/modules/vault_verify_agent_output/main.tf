# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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

variable "transport" {
  description = "The transport configuration to use when setting up the auth method. Must include the name and config, i.e. ssh = { host = ??, user = ?? }"
  type        = any # Cannot use object (even with optional properties) since it blows up.
}

resource "enos_remote_exec" "verify_vault_agent_output" {
  content = templatefile("${path.module}/templates/verify-vault-agent-output.sh", {
    vault_agent_template_destination = var.vault_agent_template_destination
    vault_agent_expected_output      = var.vault_agent_expected_output
  })

  transport = var.transport
}
