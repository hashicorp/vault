# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

// An arithmetic module for calculating inputs and outputs for various replication steps.

// Get the first follower out of the hosts set
variable "follower_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  default = {}
}

output "follower_host_1" {
  value = try(var.follower_hosts[0], null)
}

output "follower_public_ip_1" {
  value = try(var.follower_hosts[0].public_ip, null)
}

output "follower_private_ip_1" {
  value = try(var.follower_hosts[0].private_ip, null)
}

output "follower_host_2" {
  value = try(var.follower_hosts[1], null)
}

output "follower_public_ip_2" {
  value = try(var.follower_hosts[1].public_ip, null)
}

output "follower_private_ip_2" {
  value = try(var.follower_hosts[1].private_ip, null)
}

// Calculate our remainder hosts after we've added and removed leader
variable "initial_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  default = {}
}

variable "initial_hosts_count" {
  type    = number
  default = 0
}

variable "added_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  default = {}
}

variable "added_hosts_count" {
  type    = number
  default = 0
}

variable "removed_primary_host" {
  type = object({
    private_ip = string
    public_ip  = string
  })
  default = null
}

variable "removed_follower_host" {
  type = object({
    private_ip = string
    public_ip  = string
  })
  default = null
}

locals {
  remaining_hosts_count = max((var.initial_hosts_count + var.added_hosts_count - 2), 0)
  indices               = [for idx in range(local.remaining_hosts_count) : idx]
  remaining_initial     = setsubtract(values(var.initial_hosts), [var.removed_primary_host, var.removed_follower_host])
  remaining_hosts_list  = tolist(setunion(values(var.added_hosts), local.remaining_initial))
  remaining_hosts       = zipmap(local.indices, local.remaining_hosts_list)
}

output "remaining_initial_count" {
  value = length(local.remaining_initial)
}

output "remaining_initial_hosts" {
  value = local.remaining_initial
}

output "remaining_hosts_count" {
  value = local.remaining_hosts_count
}

output "remaining_hosts" {
  value = local.remaining_hosts
}
