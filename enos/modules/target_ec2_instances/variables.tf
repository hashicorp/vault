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
  default     = { "Project" : "vault-ci" }
}

variable "disable_selinux" {
  description = "Optionally disable SELinux for certain distros/versions"
  type        = bool
  default     = true
}

variable "ebs_optimized" {
  description = "Apply EBS optimization and high throughput disks to maximize IO performance"
  type        = bool
  default     = false
}

variable "instance_count" {
  description = "The number of target instances to create"
  type        = number
  default     = 3
}

variable "instance_types" {
  description = "The instance types to use depending on architecture"
  type = object({
    amd64 = string
    arm64 = string
  })
  default = {
    amd64 = "t3a.medium"
    arm64 = "t4g.medium"
  }
}

variable "ports_ingress" {
  description = "Ports mappings to allow for ingress"
  type = list(object({
    description = string
    port        = number
    protocol    = string
  }))
}

variable "project_name" {
  description = "A unique project name"
  type        = string
}

// These 3 root volume variables all default to null so that, if they're not specified, we use the AMI default
// Note that for IOPs specifically, the ratio of IOPs to size for io1 disks is 50:1 and for io2 disks it's 1000:1
variable "root_volume_iops" {
  description = "The IOPS of the root volume"
  type        = number
  default     = null
}

variable "root_volume_size" {
  description = "The size of the root volume"
  type        = number
  default     = null
}

variable "root_volume_type" {
  description = "The type of the root volume"
  type        = string
  default     = null
}

variable "seal_key_names" {
  type        = list(string)
  description = "The key management seal key names"
  default     = []
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
