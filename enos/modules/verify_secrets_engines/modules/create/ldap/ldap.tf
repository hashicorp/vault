# Copyright IBM Corp. 2016, 2025
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

# Setup OpenLDAP infrastructure 
resource "enos_remote_exec" "ldap_setup" {
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

  scripts = [abspath("${path.module}/../../../scripts/ldap/setup.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Configure LDAP secrets engine (separate from auth backend)
resource "enos_remote_exec" "ldap_secrets_config" {
  depends_on = [
    enos_remote_exec.secrets_enable_ldap_secret,
    enos_remote_exec.ldap_setup,
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

  scripts = [abspath("${path.module}/../../../scripts/ldap-secrets-config.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Create a new Library set of service accounts
# Test Case: Service Account Library - Create a new Library set of service accounts
resource "enos_remote_exec" "ldap_library_set_create" {
  depends_on = [
    enos_remote_exec.ldap_secrets_config,
  ]

  environment = {
    REQPATH = "${local.ldap_output.ldap_mount}/library/test-set"
    PAYLOAD = jsonencode({
      service_account_names        = "fizz"
      ttl                          = "10h"
      max_ttl                      = "20h"
      disable_check_in_enforcement = false
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Update Library configuration
# Test Case: Modify library settings and accounts - Update library configuration and associated service accounts
resource "enos_remote_exec" "ldap_library_set_update" {
  depends_on = [
    enos_remote_exec.ldap_library_set_create,
  ]

  environment = {
    REQPATH = "${local.ldap_output.ldap_mount}/library/test-set"
    PAYLOAD = jsonencode({
      service_account_names        = "fizz,buzz"
      ttl                          = "12h"
      max_ttl                      = "15h"
      disable_check_in_enforcement = true
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Check-out Service Account
# Test Case: Check-out Service Account - Borrow service accounts for temporary use
resource "enos_remote_exec" "ldap_library_checkout_default_ttl" {
  depends_on = [
    enos_remote_exec.ldap_library_set_update,
  ]

  environment = {
    REQPATH           = "${local.ldap_output.ldap_mount}/library/test-set/check-out"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Check-out with Custom TTL
# Test Case: Check-out with Custom TTL - Borrow with specific lease duration
resource "enos_remote_exec" "ldap_library_checkout_custom_ttl" {
  depends_on = [
    enos_remote_exec.ldap_library_checkout_default_ttl,
  ]

  environment = {
    REQPATH = "${local.ldap_output.ldap_mount}/library/test-set/check-out"
    PAYLOAD = jsonencode({
      ttl = "2h"
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Self Check-in (Explicit)
# Test Case: Self Check-in (Explicit) - Return your checked-out account
resource "enos_remote_exec" "ldap_library_self_checkin" {
  depends_on = [
    enos_remote_exec.ldap_library_checkout_custom_ttl,
  ]

  environment = {
    REQPATH = "${local.ldap_output.ldap_mount}/library/test-set/check-in"
    PAYLOAD = jsonencode({
      service_account_names = "fizz,buzz"
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

resource "enos_remote_exec" "ldap_password_policy" {
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

  scripts = [abspath("${path.module}/../../../scripts/ldap/add-ldap-password-policy.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
