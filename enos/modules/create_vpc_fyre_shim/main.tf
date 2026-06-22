# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

variable "name" {
  type        = string
  default     = "vault-ci"
  description = "The synthetic environment name"
}

variable "ip_version" {
  type        = number
  default     = 4
  description = "The IP version to use for compatibility outputs"

  validation {
    # Fyre only supports ipv4 but we want to have compatible module interfaces
    condition     = var.ip_version == 4
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "ipv4_cidr" {
  type        = string
  default     = "9.13.0.0/16"
  description = "Compatibility CIDR block for ipv4 mode"
}

variable "ipv6_cidr" {
  type        = string
  default     = ""
  description = "Compatibility CIDR block for ipv6 mode"
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

resource "random_string" "cluster_id" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

output "id" {
  description = "Synthetic environment ID for Fyre-backed scenarios"
  value       = "${var.name}-${random_string.cluster_id.result}"
}

output "ipv4_cidr" {
  description = "Compatibility IPv4 CIDR value for downstream modules"
  value       = var.ipv4_cidr
}

output "ipv6_cidr" {
  description = "Compatibility IPv6 CIDR value for downstream modules"
  value       = null
}

output "cluster_id" {
  description = "A unique string associated with the synthetic Fyre environment"
  value       = random_string.cluster_id.result
}
