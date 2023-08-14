# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "repository" {
  description = "The GitHub repository, either vault or vault-enterprise"
  type        = string
  validation {
    condition     = contains(["vault", "vault-enterprise"], var.repository)
    error_message = "Invalid repository, only vault or vault-enterprise are supported"
  }
}
