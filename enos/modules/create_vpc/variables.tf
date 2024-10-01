# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "name" {
  type        = string
  default     = "vault-ci"
  description = "The name of the VPC"
}

variable "ip_version" {
  type        = number
  default     = 4
  description = "The IP version to use for the default subnet"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "ipv4_cidr" {
  type        = string
  default     = "10.13.0.0/16"
  description = "The CIDR block for the VPC when using IPV4 mode"
}

variable "environment" {
  description = "Name of the environment."
  type        = string
  default     = "vault-ci"
}

variable "common_tags" {
  description = "Tags to set for all resources"
  type        = map(string)
  default     = { "Project" : "vault-ci" }
}
