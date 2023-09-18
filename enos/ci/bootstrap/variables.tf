# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "aws_ssh_public_key" {
  description = "The public key to use for the ssh key"
  type        = string
}

variable "repository" {
  description = "The repository to bootstrap the ci for, either 'vault' or 'vault-enterprise'"
  type        = string
  validation {
    condition     = contains(["vault", "vault-enterprise"], var.repository)
    error_message = "Repository must be one of either 'vault' or 'vault-enterprise'"
  }
}
