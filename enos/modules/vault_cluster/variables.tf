# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "artifactory_release" {
  type = object({
    username = string
    token    = string
    url      = string
    sha256   = string
  })
  description = "The Artifactory release information to install Vault artifacts from Artifactory"
  default     = null
}

variable "backend_cluster_name" {
  type        = string
  description = "The name of the backend cluster"
  default     = null
}

variable "backend_cluster_tag_key" {
  type        = string
  description = "The tag key for searching for backend nodes"
  default     = null
}

variable "cluster_name" {
  type        = string
  description = "The Vault cluster name"
  default     = null
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

variable "config_env_vars" {
  description = "Optional Vault configuration environment variables to set starting Vault"
  type        = map(string)
  default     = null
}

variable "consul_data_dir" {
  type        = string
  description = "The directory where the consul will store data"
  default     = "/opt/consul/data"
}

variable "consul_install_dir" {
  type        = string
  description = "The directory where the consul binary will be installed"
  default     = "/opt/consul/bin"
}

variable "consul_license" {
  type        = string
  sensitive   = true
  description = "The consul enterprise license"
  default     = null
}

variable "consul_log_file" {
  type        = string
  description = "The file where the consul will write log output"
  default     = "/var/log/consul.log"
}

variable "consul_log_level" {
  type        = string
  description = "The consul service log level"
  default     = "info"

  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.consul_log_level)
    error_message = "The consul_log_level must be one of 'trace', 'debug', 'info', 'warn', or 'error'."
  }
}

variable "consul_release" {
  type = object({
    version = string
    edition = string
  })
  description = "Consul release version and edition to install from releases.hashicorp.com"
  default = {
    version = "1.15.1"
    edition = "ce"
  }
}

variable "distro_version" {
  type        = string
  description = "The Linux distro version"
  default     = null
}

variable "enable_audit_devices" {
  description = "If true every audit device will be enabled"
  type        = bool
  default     = true
}

variable "force_unseal" {
  type        = bool
  description = "Always unseal the Vault cluster even if we're not initializing it"
  default     = false
}

variable "initialize_cluster" {
  type        = bool
  description = "Initialize the Vault cluster"
  default     = true
}

variable "install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = "/opt/vault/bin"
}

variable "license" {
  type        = string
  sensitive   = true
  description = "The value of the Vault license"
  default     = null
}

variable "local_artifact_path" {
  type        = string
  description = "The path to a locally built vault artifact to install. It can be a zip archive, RPM, or Debian package"
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

variable "packages" {
  type        = list(string)
  description = "A list of packages to install via the target host package manager"
  default     = []
}

variable "release" {
  type = object({
    version = string
    edition = string
  })
  description = "Vault release version and edition to install from releases.hashicorp.com"
  default     = null
}

variable "root_token" {
  type        = string
  description = "The Vault root token that we can use to intialize and configure the cluster"
  default     = null
}

variable "seal_ha_beta" {
  description = "Enable using Seal HA on clusters that meet minimum version requirements and are enterprise editions"
  default     = true
}

variable "seal_attributes" {
  description = "The auto-unseal device attributes"
  default     = null
}

variable "seal_attributes_secondary" {
  description = "The secondary auto-unseal device attributes"
  default     = null
}

variable "seal_type" {
  type        = string
  description = "The primary seal device type"
  default     = "awskms"

  validation {
    condition     = contains(["awskms", "pkcs11", "shamir"], var.seal_type)
    error_message = "The seal_type must be either 'awskms', 'pkcs11', or 'shamir'. No other seal types are supported."
  }
}

variable "seal_type_secondary" {
  type        = string
  description = "A secondary HA seal device type. Only supported in Vault Enterprise >= 1.15"
  default     = "none"

  validation {
    condition     = contains(["awskms", "none", "pkcs11"], var.seal_type_secondary)
    error_message = "The secondary_seal_type must be 'awskms', 'none', or 'pkcs11'. No other secondary seal types are supported."
  }
}

variable "shamir_unseal_keys" {
  type        = list(string)
  description = "Shamir unseal keys. Often only used adding additional nodes to an already initialized cluster."
  default     = null
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

variable "storage_backend_addl_config" {
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
