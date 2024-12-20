# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster followers"
}


variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
  default     = 60
}

variable "listener_port" {
  type        = number
  description = "The listener port for vault"
}
variable "vault_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The leader's host information"
}
variable "vault_local_addr" {
  type        = string
  description = "The local address to use to query vault"
}
variable "cluster_port" {
  type        = number
  description = "The cluster port for vault"
}

variable "ip_version" {
  type        = number
  description = "The IP version to use for the Vault TCP listeners"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}
variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}
variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "add_back_nodes" {
  type        = bool
  description = "whether to add the nodes back"
}

variable "vault_unseal_keys" {}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the vault binary is installed"
}


module "choose_follower_to_remove" {
  source    = "../choose_follower_host"
  followers = var.hosts
}

module "remove_raft_node" {
  source     = "../vault_raft_remove_peer"
  depends_on = [module.choose_follower_to_remove]


  hosts                   = module.choose_follower_to_remove.chosen_follower
  ip_version              = var.ip_version
  operator_instance       = var.vault_leader_host.public_ip
  vault_addr              = var.vault_local_addr
  vault_cluster_addr_port = var.cluster_port
  vault_install_dir       = var.vault_install_dir
  vault_root_token        = var.vault_root_token
  is_voter                = true
}

module "verify_removed" {
  source     = "../vault_verify_raft_removed"
  depends_on = [
    module.remove_raft_node
  ]

  hosts             = module.choose_follower_to_remove.chosen_follower
  vault_leader_host = var.vault_leader_host
  vault_root_token  = var.vault_root_token
  vault_seal_type   = var.vault_seal_type
  vault_unseal_keys = var.vault_seal_type == "shamir" ? var.vault_unseal_keys : null
  add_back_nodes    = true
  listener_port     = var.listener_port
  ip_version        = var.ip_version
  vault_local_addr  = var.vault_local_addr
  cluster_port      = var.cluster_port
  vault_install_dir = var.vault_install_dir
}
