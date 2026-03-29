# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# Plugin blackbox test configuration

variable "plugin_config" {
  type = object({
    enabled = bool
    type    = string
  })
  description = "Plugin configuration for blackbox tests. Set enabled=true and specify plugin type to enable plugin tests."
  default = {
    enabled = false
    type    = ""
  }
}


# Local variables for plugin environment setup
locals {
  plugin_environment = var.plugin_config.enabled ? {
    PLUGIN_TYPE = var.plugin_config.type
  } : {}
}
