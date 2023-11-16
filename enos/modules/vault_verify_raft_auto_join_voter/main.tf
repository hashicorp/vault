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

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "verify_raft_auto_join_voter" {
  for_each = local.instances

  environment = {
    VAULT_CLUSTER_ADDR      = "${each.value.private_ip}:${var.vault_cluster_addr_port}"
    VAULT_INSTALL_DIR       = var.vault_install_dir
    VAULT_LOCAL_BINARY_PATH = "${var.vault_install_dir}/vault"
    VAULT_TOKEN             = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/verify-raft-auto-join-voter.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
