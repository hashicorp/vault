// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

variable "dev_build_local_ui" {
  type        = bool
  description = "Whether or not to build the web UI when using the local builder var. If the assets have already been built we'll still include them"
  default     = false
}

variable "dev_config_mode" {
  type        = string
  description = "The method to use when configuring Vault. When set to 'env' we will configure Vault using VAULT_ style environment variables if possible. When 'file' we'll use the HCL configuration file for all configuration options."
  default     = "file" // or "env"
}

variable "dev_consul_version" {
  type        = string
  description = "The version of Consul to use when using Consul for storage!"
  default     = "1.18.1"
  // NOTE: You can also set the "backend_edition" if you want to use Consul Enterprise
}
