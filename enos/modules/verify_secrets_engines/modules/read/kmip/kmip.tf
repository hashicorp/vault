# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "create_state" {
  description = "The state of the secrets engines from the 'create' module"
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "ip_version" {
  type        = string
  description = "IP Version (4 or 6)"
  default     = "4"
}

variable "vault_edition" {
  type        = string
  description = "IP Version (4 or 6)"
  default     = "4"
}

locals {
  kmip_output = {
    mount      = "kmip"
    ip_address = var.ip_version == "6" ? var.hosts[0].ipv6 : var.hosts[0].public_ip
  }
}

# KMIP Client Configuration
resource "enos_remote_exec" "kmip_client_configure" {

  environment = {
    VAULT_ADDR = var.vault_addr
    SERVER_CA  = var.create_state.kmip.server_ca
    CLIENT_CA  = var.create_state.kmip.client_cert
    KMIP_PORT  = var.create_state.kmip.port
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../../scripts/kmip/kmip-client-configure.sh")]

  transport = {
    ssh = {
      host = var.create_state.kmip.test_server_ip
      user = "ubuntu" # Assuming Ubuntu for the test server
    }
  }
}
