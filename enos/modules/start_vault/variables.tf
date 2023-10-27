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

variable "seal_ha_beta" {
  description = "Enable using Seal HA on clusters that meet minimum version requirements and are enterprise editions"
  default     = true
}

variable "seal_key_name" {
  type        = string
  description = "The auto-unseal key name"
  default     = null
}

variable "seal_key_name_secondary" {
  type        = string
  description = "The secondary auto-unseal key name"
  default     = null
}

variable "seal_type" {
  type        = string
  description = "The method by which to unseal the Vault cluster"
  default     = "awskms"

  validation {
    condition     = contains(["awskms", "shamir"], var.seal_type)
    error_message = "The seal_type must be either awskms or shamir. No other unseal methods are supported."
  }
}

variable "seal_type_secondary" {
  type        = string
  description = "A secondary HA seal method. Only supported in Vault Enterprise >= 1.15"
  default     = "none"

  validation {
    condition     = contains(["awskms", "none"], var.seal_type_secondary)
    error_message = "The secondary_seal_type must be 'awskms' or 'none'. No other secondary unseal methods are supported."
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
