# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

// An arithmetic module for calculating inputs and outputs for various replication steps.

variable "added_hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  default = {}
}

variable "initial_hosts" {
  description = "The initial set of Vault cluster hosts before removing and adding hosts"
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  default = {}
}

variable "removed_primary_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  default = null
}

variable "removed_follower_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  default = null
}

locals {
  remaining_initial    = setsubtract(values(var.initial_hosts), [var.removed_primary_host, var.removed_follower_host])
  remaining_hosts_list = tolist(setunion(values(var.added_hosts), local.remaining_initial))
  remaining_hosts      = { for idx in range(length(local.remaining_hosts_list)) : idx => local.remaining_hosts_list[idx] }
}

output "remaining_hosts" {
  value = local.remaining_hosts
}
