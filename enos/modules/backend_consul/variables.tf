# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "cluster_name" {
  type        = string
  description = "The name of the Consul cluster"
  default     = null
}

variable "cluster_tag_key" {
  type        = string
  description = "The tag key for searching for Consul nodes"
  default     = null
}

variable "config_dir" {
  type        = string
  description = "The directory where the consul will write config files"
  default     = "/etc/consul.d"
}

variable "data_dir" {
  type        = string
  description = "The directory where the consul will store data"
  default     = "/opt/consul/data"
}

variable "install_dir" {
  type        = string
  description = "The directory where the consul binary will be installed"
  default     = "/opt/consul/bin"
}

variable "license" {
  type        = string
  sensitive   = true
  description = "The consul enterprise license"
  default     = null
}

variable "log_dir" {
  type        = string
  description = "The directory where the consul will write log files"
  default     = "/var/log/consul.d"
}

variable "log_level" {
  type        = string
  description = "The consul service log level"
  default     = "info"

  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.log_level)
    error_message = "The log_level must be one of 'trace', 'debug', 'info', 'warn', or 'error'."
  }
}

variable "release" {
  type = object({
    version = string
    edition = string
  })
  description = "Consul release version and edition to install from releases.hashicorp.com"
  default = {
    version = "1.15.3"
    edition = "oss"
  }
}

variable "target_hosts" {
  description = "The target machines host addresses to use for the consul cluster"
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
}
