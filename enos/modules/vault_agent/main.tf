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

variable "ip_version" {
  type        = number
  default     = 4
  description = "The IP version to use for the Vault TCP listeners"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_agent_port" {
  type        = number
  description = "The listener port number for the Vault Agent"
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

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

locals {
  agent_listen_addr = "${var.ip_version == 4 ? "127.0.0.1" : "[::1]"}:${var.vault_agent_port}"
}

resource "enos_remote_exec" "set_up_approle_auth_and_agent" {
  environment = {
    AGENT_LISTEN_ADDR                = local.agent_listen_addr,
    VAULT_ADDR                       = var.vault_addr,
    VAULT_INSTALL_DIR                = var.vault_install_dir,
    VAULT_TOKEN                      = var.vault_root_token,
    VAULT_AGENT_TEMPLATE_DESTINATION = var.vault_agent_template_destination,
    VAULT_AGENT_TEMPLATE_CONTENTS    = var.vault_agent_template_contents,
  }

  scripts = [abspath("${path.module}/scripts/set-up-approle-and-agent.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

output "vault_agent_listen_addr" {
  description = "The vault agent listen address"
  value       = local.agent_listen_addr
}
