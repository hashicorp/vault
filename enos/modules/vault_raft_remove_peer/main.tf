# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
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
  description = "The old vault nodes to be removed"
}

variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "operator_instance" {
  type        = string
  description = "The ip address of the operator (Voter) node"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_cluster_addr_port" {
  description = "The Raft cluster address port"
  type        = string
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "is_voter" {
  type        = bool
  default     = false
  description = "Whether the nodes that are going to be removed are voters"
}

resource "enos_remote_exec" "vault_raft_remove_peer" {
  for_each = var.hosts

  environment = {
    REMOVE_VAULT_CLUSTER_ADDR = "${var.ip_version == 4 ? "${each.value.private_ip}" : "[${each.value.ipv6}]"}:${var.vault_cluster_addr_port}"
    VAULT_TOKEN               = var.vault_root_token
    VAULT_ADDR                = var.vault_addr
    VAULT_INSTALL_DIR         = var.vault_install_dir
    REMOVE_NODE_IS_VOTER      = var.is_voter
  }

  scripts = [abspath("${path.module}/scripts/raft-remove-peer.sh")]

  transport = {
    ssh = {
      host = var.operator_instance
    }
  }
}
