# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_agent_template_destination" {
  type        = string
  description = "The destination of the template rendered by Agent"
}

variable "vault_agent_template_contents" {
  type        = string
  description = "The template contents to be rendered by Agent"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

locals {
  vault_instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "set_up_approle_auth_and_agent" {
  content = templatefile("${path.module}/templates/set-up-approle-and-agent.sh", {
    vault_install_dir                = var.vault_install_dir
    vault_token                      = var.vault_root_token
    vault_agent_template_destination = var.vault_agent_template_destination
    vault_agent_template_contents    = var.vault_agent_template_contents
  })

  transport = {
    ssh = {
      host = local.vault_instances[0].public_ip
    }
  }
}
