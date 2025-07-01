# Copyright (c) HashiCorp, Inc.
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

variable "distro" {
  type        = string
  description = "OS Distribution"
  default     = "ubuntu"
}