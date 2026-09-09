# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

variable "cluster_name" {
  type        = string
  description = "A unique cluster identifier"
  default     = null
}

variable "cluster_tag_key" {
  type        = string
  description = "A logical key name used in generated hostnames"
  default     = "TargetCluster"
}

variable "common_tags" {
  description = "Common tags for cloud resources"
  type        = map(string)
  default     = { "Project" : "vault-ci" }
}

variable "cpu" {
  description = "The number of CPUs for each Fyre VM"
  type        = number
  default     = 2
}

variable "description" {
  description = "Optional description applied to created Fyre VMs"
  type        = string
  default     = null
}

variable "disable_delete" {
  description = "Prevent accidental deletion of created Fyre VMs"
  type        = string
  default     = "n"

  validation {
    condition     = contains(["y", "n"], var.disable_delete)
    error_message = "The disable_delete value must be either 'y' or 'n'."
  }
}

variable "disable_selinux" {
  description = "Optionally disable SELinux for certain distros/versions"
  type        = bool
  default     = true
}

variable "dns" {
  description = "Whether to add the VM hostname to DNS"
  type        = string
  default     = "n"

  validation {
    condition     = contains(["y", "n"], var.dns)
    error_message = "The dns value must be either 'y' or 'n'."
  }
}

variable "additional_disks" {
  description = "Additional disk sizes in GB"
  type        = list(string)
  default     = []
}

variable "expiration" {
  description = "VM expiration in hours or Fyre-supported duration format"
  type        = string
  default     = null
}

variable "instance_count" {
  description = "The number of target instances to create"
  type        = number
  default     = 3
}

variable "memory" {
  description = "The amount of memory in GB for each Fyre VM"
  type        = number
  default     = 4
}

variable "os" {
  description = "The Fyre OS identifier to provision"
  type        = string
}

variable "arch" {
  description = "The target architecture to provision"
  type        = string

  validation {
    condition     = contains(["amd64", "s390x"], var.arch)
    error_message = "The arch must be one of 'amd64' or 's390x'."
  }
}

variable "product_group_id" {
  description = "Optional Fyre product group override"
  type        = number
  default     = null
}

variable "project_name" {
  description = "A unique project name"
  type        = string
  default     = "vault-ci"
}

variable "public_key_path" {
  description = "Path to the RFC 4716 formatted public key used for SSH access"
  type        = string
}

variable "public_network" {
  description = "Whether to assign a public IP address"
  type        = string
  default     = "y"

  validation {
    condition     = contains(["y", "n"], var.public_network)
    error_message = "The public_network value must be either 'y' or 'n'."
  }
}

variable "quota_type" {
  description = "The Fyre quota type to use"
  type        = string
  default     = "product_group"

  validation {
    condition     = contains(["product_group", "quick_burn"], var.quota_type)
    error_message = "The quota_type must be either 'product_group' or 'quick_burn'."
  }
}

variable "quick_burn_ttl" {
  description = "TTL in hours when using quick_burn quota"
  type        = string
  default     = "6"
}
