# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })

  description = "Vault cluster leader host"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
  default     = null
}

variable "ldap_password" {
  type        = string
  description = "The LDAP Server admin password"
  default     = "password1"
}

variable "integration_host_state" {
  description = "The state of the test server from the 'set_up_external_integration' module"
}

variable "ip_version" {
  type        = string
  description = "IP Version (4 or 6)"
  default     = "4"
}

variable "ports" {
  description = "Port configuration for services"
  type = map(object({
    port        = string
    description = string
  }))
}

locals {
  ldap_output = {
    ip_version = var.ip_version
    ldap_mount = "ldap"
    host       = var.integration_host_state.ldap.host
    port       = var.ports.ldap.port
    username   = "enos"
    pw         = var.ldap_password
  }
}

output "ldap" {
  value = local.ldap_output
}

# Enable LDAP secrets engine
resource "enos_remote_exec" "secrets_enable_ldap_secret" {
  environment = {
    ENGINE            = local.ldap_output.ldap_mount
    MOUNT             = local.ldap_output.ldap_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/secrets-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Configuring Openldap Server and Vault LDAP
resource "enos_remote_exec" "ldap_configurations" {
  depends_on = [
    enos_remote_exec.secrets_enable_ldap_secret
  ]

  environment = {
    MOUNT             = local.ldap_output.ldap_mount
    LDAP_SERVER       = local.ldap_output.host.private_ip
    LDAP_PORT         = local.ldap_output.port
    LDAP_USERNAME     = local.ldap_output.username
    LDAP_ADMIN_PW     = local.ldap_output.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-configs.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
