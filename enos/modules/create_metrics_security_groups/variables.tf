# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "project_name" {
  type        = string
  default     = ""
  description = "The name of the project"
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC"
}