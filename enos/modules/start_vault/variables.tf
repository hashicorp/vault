# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "cluster_name" {
  type        = string
  description = "The Vault cluster name"
}

variable "config_dir" {
  type        = string
  description = "The directory to use for Vault configuration"
  default     = "/etc/vault.d"
}

variable "config_mode" {
  description = "The method to use when configuring Vault. When set to 'env' we will configure Vault using VAULT_ style environment variables if possible. When 'file' we'll use the HCL configuration file for all configuration options."
  default     = "file"

  validation {
    condition     = contains(["env", "file"], var.config_mode)
    error_message = "The config_mode must be either 'env' or 'file'. No other configuration modes are supported."
  }
}

variable "environment" {
  description = "Optional Vault configuration environment variables to set starting Vault"
  type        = map(string)
  default     = null
}

variable "install_dir" {
  type        = string
  description = "The directory where the vault binary will be installed"
  default     = "/opt/vault/bin"
}

variable "license" {
  type        = string
  sensitive   = true
  description = "The value of the Vault license"
  default     = null
}

variable "log_level" {
  type        = string
  description = "The vault service log level"
  default     = "info"

  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.log_level)
    error_message = "The log_level must be one of 'trace', 'debug', 'info', 'warn', or 'error'."
  }
}

variable "manage_service" {
  type        = bool
  description = "Manage the Vault service users and systemd unit. Disable this to use configuration in RPM and Debian packages"
  default     = true
}

variable "seal_alias" {
  type        = string
  description = "The primary seal alias name"
  default     = "primary"
}

variable "seal_alias_secondary" {
  type        = string
  description = "The secondary seal alias name"
  default     = "secondary"
}

variable "seal_attributes" {
  description = "The primary auto-unseal attributes"
  default     = null
}

variable "seal_attributes_secondary" {
  description = "The secondary auto-unseal attributes"
  default     = null
}

variable "seal_priority" {
  type        = string
  description = "The primary seal priority"
  default     = "1"
}

variable "seal_priority_secondary" {
  type        = string
  description = "The secondary seal priority"
  default     = "2"
}

variable "seal_type" {
  type        = string
  description = "The method by which to unseal the Vault cluster"
  default     = "awskms"

  validation {
    condition     = contains(["awskms", "pkcs11", "shamir"], var.seal_type)
    error_message = "The seal_type must be either 'awskms', 'pkcs11', or 'shamir'. No other seal types are supported."
  }
}

variable "seal_type_secondary" {
  type        = string
  description = "A secondary HA seal method. Only supported in Vault Enterprise >= 1.15"
  default     = "none"

  validation {
    condition     = contains(["awskms", "pkcs11", "none"], var.seal_type_secondary)
    error_message = "The secondary_seal_type must be 'awskms', 'pkcs11' or 'none'. No other secondary seal types are supported."
  }
}

variable "service_username" {
  type        = string
  description = "The host username to own the vault service"
  default     = "vault"
}

variable "storage_backend" {
  type        = string
  description = "The storage backend to use"
  default     = "raft"

  validation {
    condition     = contains(["raft", "consul"], var.storage_backend)
    error_message = "The storage_backend must be either raft or consul. No other storage backends are supported."
  }
}

variable "storage_backend_attrs" {
  type        = map(any)
  description = "An optional set of key value pairs to inject into the storage block"
  default     = {}
}

variable "storage_node_prefix" {
  type        = string
  description = "A prefix to use for each node in the Vault storage configuration"
  default     = "node"
}

variable "target_hosts" {
  description = "The target machines host addresses to use for the Vault cluster"
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
}
