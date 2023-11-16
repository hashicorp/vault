# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_cluster_addr_port" {
  description = "The Raft cluster address port"
  type        = string
  default     = "8201"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "operator_instance" {
  type        = string
  description = "The ip address of the operator (Voter) node"
}

variable "remove_vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The old vault nodes to be removed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.remove_vault_instances)[idx].public_ip
      private_ip = values(var.remove_vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "vault_raft_remove_peer" {
  for_each = local.instances

  environment = {
    REMOVE_VAULT_CLUSTER_ADDR = "${each.value.private_ip}:${var.vault_cluster_addr_port}"
    VAULT_TOKEN               = var.vault_root_token
    VAULT_ADDR                = "http://localhost:8200"
    VAULT_INSTALL_DIR         = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/raft-remove-peer.sh")]

  transport = {
    ssh = {
      host = var.operator_instance
    }
  }
}
