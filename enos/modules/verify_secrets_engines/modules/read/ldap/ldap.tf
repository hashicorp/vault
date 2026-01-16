# Copyright IBM Corp. 2016, 2025
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
  description = "The Vault cluster instances that were created"
}

variable "create_state" {
  description = "The state of the secrets engines from the 'create' module"
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

variable "credential_ttl_buffer" {
  description = "Buffer (seconds) to wait after LDAP credential TTL expiry"
  type        = number
  default     = 80
}

variable "default_ttl" {
  description = "Default time-to-live (in seconds) for issued LDAP credentials."
  type        = number
  default     = 60
}

variable "max_ttl" {
  description = "Maximum time-to-live (in seconds) allowed for issued LDAP credentials"
  type        = number
  default     = 60
}

variable "enable_secrets_verification" {
  type        = bool
  description = "Enable LDAP secrets engine verification (dynamic credentials)"
  default     = true
}

variable "enable_rotation_verification" {
  type        = bool
  description = "Enable LDAP root rotation verification"
  default     = true
}

variable "enable_auth_verification" {
  type        = bool
  description = "Enable LDAP authentication verification"
  default     = true
}

resource "enos_remote_exec" "ldap_verify_auth" {
  count = var.enable_auth_verification ? 1 : 0
  environment = {
    MOUNT             = "${var.create_state.ldap.ldap_mount}"
    LDAP_SERVER       = "${var.create_state.ldap.host.private_ip}"
    LDAP_PORT         = "${var.create_state.ldap.port}"
    LDAP_USERNAME     = "${var.create_state.ldap.username}"
    LDAP_ADMIN_PW     = "${var.create_state.ldap.pw}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-auth.sh")]
  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Configure and verify LDAP secrets engine 
resource "enos_remote_exec" "ldap_verify_secrets" {
  count = var.enable_secrets_verification ? 1 : 0

  environment = {
    MOUNT                 = "${var.create_state.ldap.ldap_mount}"
    LDAP_SERVER           = "${var.create_state.ldap.host.private_ip}"
    LDAP_PORT             = "${var.create_state.ldap.port}"
    LDAP_USERNAME         = "${var.create_state.ldap.username}"
    LDAP_ADMIN_PW         = "${var.create_state.ldap.pw}"
    VAULT_ADDR            = var.vault_addr
    VAULT_INSTALL_DIR     = var.vault_install_dir
    VAULT_TOKEN           = var.vault_root_token
    CREDENTIAL_TTL_BUFFER = tostring(var.credential_ttl_buffer)
    DEFAULT_TTL           = tostring(var.default_ttl)
    MAX_TTL               = tostring(var.max_ttl)
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-secrets.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Verify LDAP root rotation 
resource "enos_remote_exec" "ldap_verify_rotation" {
  count = var.enable_rotation_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets
  ]
  environment = {
    MOUNT             = "${var.create_state.ldap.ldap_mount}"
    LDAP_SERVER       = "${var.create_state.ldap.host.private_ip}"
    LDAP_PORT         = "${var.create_state.ldap.port}"
    LDAP_USERNAME     = "${var.create_state.ldap.username}"
    LDAP_ADMIN_PW     = "${var.create_state.ldap.pw}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-rotation.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

variable "enable_rollback_verification" {
  type        = bool
  description = "Enable LDAP secrets engine rollback verification"
  default     = true
}

# Configure and verify LDAP secrets engine rollback behavior
resource "enos_remote_exec" "ldap_verify_rollback" {
  count = var.enable_rollback_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_verify_rotation
  ]
  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [
    abspath("${path.module}/../../../scripts/ldap/secrets-rollback/secrets-rollback-invalid-config.sh"),
    abspath("${path.module}/../../../scripts/ldap/secrets-rollback/secrets-rollback-creds-mismatch.sh"),
    abspath("${path.module}/../../../scripts/ldap/secrets-rollback/secrets-rollback-transactional.sh"),
  ]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}


