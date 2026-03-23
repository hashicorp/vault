# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

variable "hosts" {
  description = "The target machines host addresses to use for the Vault cluster"
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
}

variable "ip_version" {
  type        = string
  description = "IP Version (4 or 6)"
  default     = "4"
}

variable "ldap_version" {
  type        = string
  description = "OpenLDAP Server Version to use"
  default     = "1.5.0"
}

variable "packages" {
  type        = list(string)
  description = "A list of packages to install via the target host package manager"
  default     = []
}

variable "ports" {
  description = "Port configuration for services"
  type = map(object({
    port        = string
    description = string
  }))
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
  default     = 60
}
