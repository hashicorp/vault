# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "ami_id" {
  description = "The machine image identifier"
  type        = string
}

variable "cluster_name" {
  type        = string
  description = "A unique cluster identifier"
  default     = null
}

variable "cluster_tag_key" {
  type        = string
  description = "The key name for the cluster tag"
  default     = "TargetCluster"
}

variable "common_tags" {
  description = "Common tags for cloud resources"
  type        = map(string)
  default = {
    Project = "Vault"
  }
}

variable "disable_selinux" {
  description = "Optionally disable SELinux for certain distros/versions"
  type        = bool
  default     = true
}

variable "instance_mem_min" {
  description = "The minimum amount of memory in mebibytes for each instance in the fleet. (1 MiB = 1024 bytes)"
  type        = number
  default     = 4096 // ~4 GB
}

variable "instance_mem_max" {
  description = "The maximum amount of memory in mebibytes for each instance in the fleet. (1 MiB = 1024 bytes)"
  type        = number
  default     = 16385 // ~16 GB
}

variable "instance_cpu_min" {
  description = "The minimum number of vCPU's for each instance in the fleet"
  type        = number
  default     = 2
}

variable "instance_cpu_max" {
  description = "The maximum number of vCPU's for each instance in the fleet"
  type        = number
  default     = 8 // Unlikely we'll ever get that high due to spot price bid protection
}

variable "instance_count" {
  description = "The number of target instances to create"
  type        = number
  default     = 3
}

variable "project_name" {
  description = "A unique project name"
  type        = string
}

variable "max_price" {
  description = "The maximum hourly price to pay for each target instance"
  type        = string
  default     = "0.0416"
}

variable "seal_key_names" {
  type        = list(string)
  description = "The key management seal key names"
  default     = null
}

variable "ssh_allow_ips" {
  description = "Allowlisted IP addresses for SSH access to target nodes. The IP address of the machine running Enos will automatically allowlisted"
  type        = list(string)
  default     = []
}

variable "ssh_keypair" {
  description = "SSH keypair used to connect to EC2 instances"
  type        = string
}

variable "vpc_id" {
  description = "The identifier of the VPC where the target instances will be created"
  type        = string
}
