# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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

variable "ssh_keypair" {
  description = "SSH keypair used to connect to EC2 instances"
  type        = string
}

variable "vpc_id" {
  description = "The identifier of the VPC where the target instances will be created"
  type        = string
}

variable "vpc_security_group_ids" {
  description = "The identifier of the VPC Security Group IDs where the server instance will be created"
  type        = string
}
