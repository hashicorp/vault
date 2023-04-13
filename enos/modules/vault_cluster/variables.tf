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

variable "awskms_unseal_key_arn" {
  type        = string
  description = "The AWSKMS key ARN if using the awskms unseal method"
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

variable "config_env_vars" {
  description = "Optional Vault configuration environment variables to set starting Vault"
  type        = map(string)
  default     = null
}

variable "consul_cluster_tag" {
  type        = string
  description = "The retry_join tag to use for Consul"
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

variable "consul_log_file" {
  type        = string
  description = "The file where the consul will write log output"
  default     = "/var/log/consul.log"
}

variable "consul_release" {
  type = object({
    version = string
    edition = string
  })
  description = "Consul release version and edition to install from releases.hashicorp.com"
  default = {
    version = "1.15.1"
    edition = "oss"
  }
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
  description = "The directory where the vault binary will be installed"
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

variable "unseal_method" {
  type        = string
  description = "The method by which to unseal the Vault cluster"
  default     = "awskms"

  validation {
    condition     = contains(["awskms", "shamir"], var.unseal_method)
    error_message = "The unseal_method must be either awskms or shamir. No other unseal methods are supported."
  }
}
